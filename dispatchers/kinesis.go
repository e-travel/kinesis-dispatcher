package dispatchers

// this is a stub struct for the moment

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
)

const KinesisBufferSize = 2 * KinesisMaxNumberOfRecords

type Kinesis struct {
	service           kinesisiface.KinesisAPI
	streamName        string
	maxBatchFrequency time.Duration
	messageQueue      chan []byte
	batchQueue        chan *KinesisBatch
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
		batchQueue:        make(chan *KinesisBatch, KinesisBufferSize),
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
	go dispatcher.processMessageQueue(NewKinesisBatch(dispatcher.streamName))
	go dispatcher.processBatchQueue()
}

func (dispatcher *Kinesis) processMessageQueue(batch *KinesisBatch) {
	ticker := time.NewTicker(dispatcher.maxBatchFrequency)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !batch.IsEmpty() {
				fmt.Println("ticker")
				dispatcher.batchQueue <- batch
				batch = NewKinesisBatch(dispatcher.streamName)
			}
		case message := <-dispatcher.messageQueue:
			fmt.Println("message")
			if batch.IsReady(message) {
				dispatcher.batchQueue <- batch
				batch = NewKinesisBatch(dispatcher.streamName)
			}
			// try to add the message
			err := batch.Add(message)
			if err != nil {
				switch err.Error() {
				case ErrKinesisBatchTooLarge:
					log.Warn("Dropping message (batch too large).")
				case ErrKinesisMessageTooLarge:
					log.Warn("Dropping message (message too large).")
				default:
					log.Error(err.Error())
				}
			}
		}
	}
}

func (dispatcher *Kinesis) processBatchQueue() {
	for batch := range dispatcher.batchQueue {
		err := dispatcher.send(batch)
		if err != nil {
			log.Error(err.Error())
		}
	}
}

func (dispatcher *Kinesis) send(batch *KinesisBatch) error {
	output, err := dispatcher.service.PutRecords(batch.records)
	if err != nil {
		return err
	}
	if *output.FailedRecordCount > 0 {
		return errors.New(fmt.Sprintf("AWS Kinesis: failed records %d/%d",
			*output.FailedRecordCount, len(batch.records.Records)))
	}
	return nil
}
