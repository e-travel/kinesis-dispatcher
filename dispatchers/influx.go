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

func sendToInflux()
	lines := bytes.Buffer{}
	lines.WriteString("cpu_load_short,host=server02 value=0.67")
	lines.WriteString("\n")
	lines.WriteString("cpu_load_short,host=server02,region=us-west value=0.55 1422568543702900257")
	lines.WriteString("\n")
	lines.WriteString("cpu_load_short,direction=in,host=server01,region=us-west value=2.0 1422568543702900257")
	lines.WriteString("\n")
	uri := "http://0.0.0.0:32772/write?db=mydb"
	resp, err := http.Post(uri, "application/x-www-form-urlencoded", &lines)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(resp.Status)

}
