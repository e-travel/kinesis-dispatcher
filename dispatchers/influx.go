package dispatchers

import "bytes"

// implements influxdb's line protocol

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
