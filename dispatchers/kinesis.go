package dispatchers

// this is a stub struct for the moment

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
)

const MEGABYTE = 1024 * 1024

// TODO: make these configurable
const awsRegion string = "eu-west-1"

const KinesisMaxNumberOfRecords = 500
const KinesisMaxSizeInBytes = 5 * MEGABYTE
const KinesisBufferSize = 2 * KinesisMaxNumberOfRecords

type Kinesis struct {
	service      kinesisiface.KinesisAPI
	streamName   string
	messageQueue chan []byte
	batchQueue   chan *kinesis.PutRecordsInput
}

func NewKinesis(streamName string) *Kinesis {
	// create session
	sess := session.Must(session.NewSession(&aws.Config{
		Retryer: client.DefaultRetryer{NumMaxRetries: 10},
		Region:  aws.String(awsRegion),
	}))
	return &Kinesis{
		service:      kinesis.New(sess),
		streamName:   streamName,
		messageQueue: make(chan []byte, KinesisBufferSize),
		batchQueue:   make(chan *kinesis.PutRecordsInput, KinesisBufferSize),
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
	messageIndex := 0
	byteCount := 0

	batch := newBatch(dispatcher.streamName)

	for message := range dispatcher.messageQueue {
		// is the batch ready?
		if isBatchReady(len(message), messageIndex+1, byteCount) {
			// enqueue the batch without blocking
			select {
			case dispatcher.batchQueue <- batch:
			default:
			}
			// reset batch
			batch = newBatch(dispatcher.streamName)
			// reset counters
			messageIndex = 0
			byteCount = 0
		}
		entry := &kinesis.PutRecordsRequestEntry{
			Data:         message,
			PartitionKey: aws.String("TODO: Change me"),
		}
		batch.Records[messageIndex] = entry
		// update counters
		byteCount += len(entry.Data) + len([]byte(*entry.PartitionKey))
		messageIndex++
	}
}

func (dispatcher *Kinesis) processBatchQueue() {
	for batch := range dispatcher.batchQueue {
		fmt.Println("We shouldn't be here")
		if output, err := dispatcher.service.PutRecords(batch); err != nil {
			fmt.Printf("error when posting to kinesis: %s\n", err.Error())
			if *output.FailedRecordCount > 0 {
				fmt.Printf("AWS Kinesis: failed records %d/%d",
					*output.FailedRecordCount, len(batch.Records))
			}
		}
	}
}

func newBatch(streamName string) *kinesis.PutRecordsInput {
	return &kinesis.PutRecordsInput{
		Records:    make([]*kinesis.PutRecordsRequestEntry, KinesisMaxNumberOfRecords),
		StreamName: aws.String(streamName),
	}
}

func isBatchReady(messageLength int, recordsLength int, byteCount int) bool {
	return byteCount+messageLength >= KinesisMaxSizeInBytes ||
		recordsLength == KinesisMaxNumberOfRecords
}
