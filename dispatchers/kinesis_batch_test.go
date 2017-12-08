package dispatchers

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/stretchr/testify/assert"
)

func TestBatch_Add_WillReturnError_WhenMessageTooLarge(t *testing.T) {
	batch := NewKinesisBatch("stream_name")
	message := strings.Repeat("a", 1+KinesisMaxSizeInBytes)
	err := batch.Add([]byte(message))
	assert.NotNil(t, err)
	assert.Equal(t, ErrKinesisMessageTooLarge, err.Error())
}

func TestBatch_Add_WillReturnError_WhenBatchTooLarge(t *testing.T) {
	batch := NewKinesisBatch("stream_name")
	message := "hello"
	for i := 0; i < KinesisMaxNumberOfRecords; i++ {
		err := batch.Add([]byte(message))
		assert.Nil(t, err)
	}
	// add one more
	err := batch.Add([]byte(message))
	assert.NotNil(t, err)
	assert.Equal(t, ErrKinesisBatchTooLarge, err.Error())
}

func TestBatch_Add_WillAppendTheMessage_AndReturnNil(t *testing.T) {
	batch := NewKinesisBatch("stream_name")
	message := "hello"
	key := generatePartitionKey([]byte(message))
	err := batch.Add([]byte(message))
	assert.Nil(t, err)
	assert.Equal(t, len(message)+len(key), batch.byteCount)
	assert.Equal(t, 1, len(batch.records.Records))
}

func TestBatch_IsReady(t *testing.T) {
	var testCases = []struct {
		rl    int // records length
		bc    int // byte count
		ml    int // message length
		ready bool
	}{
		{0, 0, 10, false},
		// no difference
		{KinesisMaxNumberOfRecords - 1, KinesisMaxSizeInBytes - 1, 0, false},
		// max number of records makes the difference
		{KinesisMaxNumberOfRecords, KinesisMaxSizeInBytes - 1, 0, true},
		// max size in bytes make the difference (due to message length)
		{KinesisMaxNumberOfRecords - 1, KinesisMaxSizeInBytes - 1, 1, true},
		// both max number of records and max size in byte make the difference
		{KinesisMaxNumberOfRecords, KinesisMaxSizeInBytes, 1, true},
	}

	data := []byte("hello")

	for _, testCase := range testCases {
		// create the batch with rl records
		batch := NewKinesisBatch("stream_name")
		for i := 0; i < testCase.rl; i++ {
			entry := &kinesis.PutRecordsRequestEntry{
				Data:         data,
				PartitionKey: aws.String(generatePartitionKey(data)),
			}
			batch.records.Records = append(batch.records.Records, entry)
		}
		// set the batch's byteCount
		batch.byteCount = testCase.bc
		// create the message with ml bytes
		message := strings.Repeat("a", testCase.ml)
		actuallyReady := batch.IsReady([]byte(message))
		assert.Equal(t, testCase.ready, actuallyReady)
	}
}

func TestBatch_IsEmpty_ReturnsTrue_WhenRecordsAreEmpty(t *testing.T) {
	batch := NewKinesisBatch("stream_name")
	assert.True(t, batch.IsEmpty())
}

func TestBatch_IsEmpty_ReturnsFalse_WhenRecordsAreNotEmpty(t *testing.T) {
	batch := NewKinesisBatch("stream_name")
	batch.Add([]byte("hello"))
	assert.False(t, batch.IsEmpty())
}

func TestBatch_Len_WillReturnTheLengthOfRecords(t *testing.T) {
	batch := NewKinesisBatch("stream_name")
	batch.Add([]byte("hello"))
	batch.Add([]byte("there"))
	assert.Equal(t, 2, batch.Len())
}

func Test_generatePartitionKey(t *testing.T) {
	var testCases = []struct {
		message []byte
		key     string
	}{
		{
			[]byte("/{-_α"),
			"/{-_α",
		},
		{
			[]byte(strings.Repeat("a", KinesisPartitionKeyMaxSize)),
			strings.Repeat("a", KinesisPartitionKeyMaxSize),
		},
		{
			[]byte(strings.Repeat("a", KinesisPartitionKeyMaxSize+1)),
			strings.Repeat("a", KinesisPartitionKeyMaxSize),
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.key, generatePartitionKey(testCase.message))
	}
}
