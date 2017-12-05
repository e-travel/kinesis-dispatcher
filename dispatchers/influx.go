package dispatchers

import (
	"io"
	"time"

	log "github.com/sirupsen/logrus"

	"bytes"
)

const InfluxMaxBatchSize = 500

// TODO: Rethink these values and make some configurable
const InfluxBatchQueueSize = 10

type Influx struct {
	client            InfluxHttpClientInterface
	batchSize         int
	maxBatchFrequency time.Duration
	messageQueue      chan []byte
	batchQueue        chan *bytes.Buffer
}

func NewInflux(client InfluxHttpClientInterface, maxBatchFrequency time.Duration) *Influx {
	return &Influx{
		client:            client,
		batchSize:         InfluxMaxBatchSize,
		maxBatchFrequency: maxBatchFrequency,
		messageQueue:      make(chan []byte, InfluxMaxBatchSize),
		batchQueue:        make(chan *bytes.Buffer, InfluxBatchQueueSize),
	}
}

func (dispatcher *Influx) Put(message []byte) bool {
	select {
	case dispatcher.messageQueue <- message:
		return true
	default:
		return false
	}
}

func (dispatcher *Influx) Dispatch() {
	go dispatcher.processMessageQueue()
	go dispatcher.processBatchQueue()
}

func (dispatcher *Influx) dispatchBatch(lines *bytes.Buffer) {
	if lines.Len() == 0 {
		return
	}
	batch := bytes.Buffer{}
	_, err := io.Copy(&batch, lines)
	if err != nil {
		log.Error(err.Error())
	} else {
		select {
		case dispatcher.batchQueue <- &batch:
		default:
			log.Error("Failed to enqueue batch")
		}
	}
	lines.Reset()
}

func (dispatcher *Influx) processMessageQueue() {
	var lines bytes.Buffer
	var lineCount int
	ticker := time.NewTicker(dispatcher.maxBatchFrequency)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if lines.Len() > 0 {
				dispatcher.dispatchBatch(&lines)
				lineCount = 0
			}
		case message := <-dispatcher.messageQueue:
			if lineCount == dispatcher.batchSize {
				dispatcher.dispatchBatch(&lines)
				lineCount = 0
			}
			// add message to buffer
			if len(message) > 0 {
				lines.Write(message)
				lines.WriteString("\n")
				lineCount++
			}
		}
	}
}

func (dispatcher *Influx) processBatchQueue() {
	for batch := range dispatcher.batchQueue {
		err := dispatcher.client.WriteLines(batch)
		if err != nil {
			log.Error(err.Error())
		}
	}
}
