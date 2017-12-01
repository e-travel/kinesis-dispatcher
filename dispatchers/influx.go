package dispatchers

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"bytes"
	"fmt"
	"net/http"
)

// implements influxdb's line protocol

// TODO: Rethink these values
const InfluxMaxBatchSize = 500
const InfluxBatchQueueSize = 10

type Influx struct {
	influxHost      string
	influxDatabase  string
	influxBatchSize int
	messageQueue    chan []byte
	batchQueue      chan bytes.Buffer
}

func NewInflux(influxHost string, influxDatabase string) *Influx {
	return &Influx{
		influxHost:      influxHost,
		influxDatabase:  influxDatabase,
		influxBatchSize: InfluxMaxBatchSize,
		messageQueue:    make(chan []byte, InfluxMaxBatchSize),
		batchQueue:      make(chan bytes.Buffer, InfluxBatchQueueSize),
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
}

func (dispatcher *Influx) processMessageQueue() {
	lines := bytes.Buffer{}
	lineCount := 0
	for message := range dispatcher.messageQueue {
		if lineCount == dispatcher.influxBatchSize {
			// send batch to influx
			select {
			case dispatcher.batchQueue <- lines:
				lines.Reset()
			default:
			}
		} else {
			// add line to buffer
			lines.Write(message)
			lines.WriteString("\n")
		}
	}
}

func (dispatcher *Influx) processBatchQueue() {
	for batch := range dispatcher.batchQueue {
		err := dispatcher.sendToInflux(batch)
		if err != nil {
			log.Error("Influx response: %s", err.Error())
		}
	}
}

func (dispatcher *Influx) sendToInflux(lines bytes.Buffer) error {
	// TODO: url escape the strings?
	uri := fmt.Sprintf("%s/write?db=%s", dispatcher.influxHost,
		dispatcher.influxDatabase)
	resp, err := http.Post(uri, "application/x-www-form-urlencoded", &lines)
	if err != nil {
		return err
	}
	// TODO: handle status properly
	switch resp.StatusCode {
	case 204:
		return nil
	case 404:
		return errors.New("Not Found")
	case 500:
		return errors.New("Internal Server Error")
	default:
		return errors.New("")
	}
}
