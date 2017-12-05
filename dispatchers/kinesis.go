package dispatchers

// this is a stub struct for the moment

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
)

const MEGABYTE = 1024 * 1024

const KinesisMaxNumberOfRecords = 500
const KinesisMaxSizeInBytes = 5 * MEGABYTE
const KinesisBufferSize = 2 * KinesisMaxNumberOfRecords
const KinesisPartitionKeyMaxSize = 256

type Kinesis struct {
	service           kinesisiface.KinesisAPI
	streamName        string
	maxBatchFrequency time.Duration
	messageQueue      chan []byte
	batchQueue        chan *kinesis.PutRecordsInput
}

func NewKinesis(streamName string, awsRegion string, maxBatchFrequency time.Duration) *Kinesis {
	// create session
	sess := session.Must(session.NewSession(&aws.Config{
		Retryer: client.DefaultRetryer{NumMaxRetries: 10},
		Region:  aws.String(awsRegion),
	}))
	return &Kinesis{
		service:           kinesis.New(sess),
		streamName:        streamName,
		maxBatchFrequency: maxBatchFrequency,
		messageQueue:      make(chan []byte, KinesisBufferSize),
		batchQueue:        make(chan *kinesis.PutRecordsInput, KinesisBufferSize),
	}
}

func (dispatcher *Kinesis) Put(message []byte) bool {
	select {
	case dispatcher.messageQueue <- message:
		return true
	default:
		return false
	}
}

func (dispatcher *Kinesis) Dispatch() {
	go dispatcher.processMessageQueue()
	go dispatcher.processBatchQueue()
}

func (dispatcher *Kinesis) processMessageQueue() {
	batch := newBatch(dispatcher.streamName)
	byteCount := 0

	ticker := time.NewTicker(dispatcher.maxBatchFrequency)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if len(batch.Records) > 0 {
				dispatcher.batchQueue <- batch
				batch = newBatch(dispatcher.streamName)
				byteCount = 0
			}
		case message := <-dispatcher.messageQueue:
			if isBatchReady(len(batch.Records), len(message), byteCount) {
				select {
				case dispatcher.batchQueue <- batch:
				default:
				}
				batch = newBatch(dispatcher.streamName)
				byteCount = 0
			}
			entry := &kinesis.PutRecordsRequestEntry{
				Data:         message,
				PartitionKey: aws.String(generatePartitionKey(message)),
			}
			// TODO: validate that the individual entry size is not above 1MB
			batch.Records = append(batch.Records, entry)
			byteCount += len(entry.Data) + len([]byte(*entry.PartitionKey))
		}
	}
}

func (dispatcher *Kinesis) processBatchQueue() {
	for batch := range dispatcher.batchQueue {
		if output, err := dispatcher.service.PutRecords(batch); err != nil {
			log.Error("error when posting to kinesis: %s\n", err.Error())
		} else {
			// TODO: for log.Debug
			// for _, record := range output.Records {
			// 	fmt.Println(record)
			// }
			if *output.FailedRecordCount > 0 {
				log.Warning("AWS Kinesis: failed records %d/%d",
					*output.FailedRecordCount, len(batch.Records))
			}
		}
	}
}

func newBatch(streamName string) *kinesis.PutRecordsInput {
	return &kinesis.PutRecordsInput{
		Records:    make([]*kinesis.PutRecordsRequestEntry, 0, KinesisMaxNumberOfRecords),
		StreamName: aws.String(streamName),
	}
}

func isBatchReady(recordsLength int, messageLength int, byteCount int) bool {
	// TODO: add some timer to the condition
	return recordsLength == KinesisMaxNumberOfRecords ||
		byteCount+messageLength >= KinesisMaxSizeInBytes
}

func generatePartitionKey(message []byte) string {
	r := []rune(string(message))
	if len(r) > KinesisPartitionKeyMaxSize {
		r = r[:KinesisPartitionKeyMaxSize]
	}
	return string(r)
}
