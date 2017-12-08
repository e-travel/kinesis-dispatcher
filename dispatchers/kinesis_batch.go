package dispatchers

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

const MEGABYTE = 1024 * 1024

const KinesisMaxNumberOfRecords = 500
const KinesisMaxSizeInBytes = 5 * MEGABYTE
const KinesisPartitionKeyMaxSize = 256

type KinesisBatch struct {
	records   *kinesis.PutRecordsInput
	byteCount int
}

func NewKinesisBatch(streamName string) *KinesisBatch {
	return &KinesisBatch{
		records: &kinesis.PutRecordsInput{
			Records:    make([]*kinesis.PutRecordsRequestEntry, 0, KinesisMaxNumberOfRecords),
			StreamName: aws.String(streamName),
		},
	}
}

// inserts message into batch; if not possible returns an error
func (batch *KinesisBatch) Add(message []byte) error {
	if len(message) > MEGABYTE {
		return errors.New(ErrKinesisMessageTooLarge)
	}
	if batch.IsReady(message) {
		return errors.New(ErrKinesisBatchTooLarge)
	}
	entry := &kinesis.PutRecordsRequestEntry{
		Data:         message,
		PartitionKey: aws.String(generatePartitionKey(message)),
	}
	batch.records.Records = append(batch.records.Records, entry)
	batch.byteCount += len(entry.Data) + len([]byte(*entry.PartitionKey))
	return nil
}

func (batch *KinesisBatch) IsReady(message []byte) bool {
	return batch.Len() == KinesisMaxNumberOfRecords ||
		batch.byteCount+len(message) >= KinesisMaxSizeInBytes
}

func (batch *KinesisBatch) IsEmpty() bool {
	return batch.Len() == 0
}

func (batch *KinesisBatch) Len() int {
	return len(batch.records.Records)
}

// TODO: revisit this function
func generatePartitionKey(message []byte) string {
	r := []rune(string(message))
	if len(r) > KinesisPartitionKeyMaxSize {
		r = r[:KinesisPartitionKeyMaxSize]
	}
	return string(r)
}
