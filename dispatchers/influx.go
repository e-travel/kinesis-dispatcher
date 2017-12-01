package dispatchers

import (
	"io"

	log "github.com/sirupsen/logrus"

	"bytes"
	"fmt"
)

// implements influxdb's line protocol
type influxClient struct {
}

// TODO: Rethink these values
const InfluxMaxBatchSize = 500
const InfluxBatchQueueSize = 10

type Influx struct {
	client          InfluxHttpClientInterface
	influxBatchSize int
	messageQueue    chan []byte
	batchQueue      chan *bytes.Buffer
}

func NewInflux(client InfluxHttpClientInterface) *Influx {
	return &Influx{
		client:          client,
		influxBatchSize: InfluxMaxBatchSize,
		messageQueue:    make(chan []byte, InfluxMaxBatchSize),
		batchQueue:      make(chan *bytes.Buffer, InfluxBatchQueueSize),
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

func (dispatcher *Influx) processMessageQueue() {
	var lines bytes.Buffer
	lineCount := 0
	for message := range dispatcher.messageQueue {
		fmt.Println("Received message")
		if lineCount == dispatcher.influxBatchSize {
			fmt.Println("Creating Batch")
			batch := bytes.Buffer{}
			_, err := io.Copy(&batch, &lines)
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
			lineCount = 0
		}
		// add line to buffer
		lines.Write(message)
		lines.WriteString("\n")
		lineCount++
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
