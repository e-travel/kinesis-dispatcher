package dispatchers

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func sendMessage(dispatcher *Kinesis, messageContent string) {
	dispatcher.Put([]byte(messageContent))
}

func fillMessageBuffer(dispatcher *Kinesis, messageContent string) {
	for i := 0; i < KinesisMaxNumberOfRecords; i++ {
		sendMessage(dispatcher, messageContent)
	}
}

func TestKinesisConstants(t *testing.T) {
	assert.Equal(t, 500, KinesisMaxNumberOfRecords)
	assert.Equal(t, 5*MEGABYTE, KinesisMaxSizeInBytes)
	assert.Equal(t, 1000, KinesisBufferSize)
}

func TestKinesis_Put_PlacesMessageToQueue(t *testing.T) {
	dispatcher := NewKinesis("stream_name", "region", 10*time.Second)
	dispatcher.messageQueue = make(chan []byte, 1)
	assert.True(t, dispatcher.Put([]byte("hello")))
	assert.Equal(t, []byte("hello"), <-dispatcher.messageQueue)
}

func TestKinesis_Put_DropsMessage_WhenQueueIsFull(t *testing.T) {
	dispatcher := NewKinesis("stream_name", "region", 10*time.Second)
	dispatcher.messageQueue = make(chan []byte, 1)
	dispatcher.Put([]byte("hello"))
	assert.False(t, dispatcher.Put([]byte("goodbye")))
	assert.Equal(t, []byte("hello"), <-dispatcher.messageQueue)
}

func TestKinesis_Dispatch_WillProcessAllQueues(t *testing.T) {
	dispatcher := NewKinesis("stream_name", "region", 10*time.Second)
	fillMessageBuffer(dispatcher, "hello")
	sendMessage(dispatcher, "hello")
	go dispatcher.Dispatch()
	// drain sink
	timer := time.NewTimer(time.Second)
	select {
	case <-dispatcher.batchQueue:
	case <-timer.C:
		assert.Fail(t, "Timer expired")
	}
	assert.Empty(t, dispatcher.messageQueue)
	assert.Empty(t, dispatcher.batchQueue)
}

func TestKinesis_ProcessMessageQueue_WillPutInBatchQueue_WhenReady(t *testing.T) {
	//t.Skip("This blocks; needs fixing")
	dispatcher := NewKinesis("stream_name", "region", 10*time.Second)
	go dispatcher.processMessageQueue()
	// create a batch by filling the buffer
	fillMessageBuffer(dispatcher, "The same message all over again")
	// send one more message to trigger batch creation
	sendMessage(dispatcher, "This will stay in the queue")
	// get the batch
	batch := <-dispatcher.batchQueue
	assert.Equal(t, KinesisMaxNumberOfRecords, len(batch.Records))
}

func TestKinesis_ProcessMessageQueue_WillPutInBatchQueue_WhenTimerFires(t *testing.T) {
	t.Skip("Implement timer test")
	dispatcher := NewKinesis("stream_name", "region", time.Microsecond)
	go dispatcher.processMessageQueue()
	// send one message to trigger batch creation
	sendMessage(dispatcher, "This will stay in the queue")
	// get the batch
	batch := <-dispatcher.batchQueue
	assert.Equal(t, KinesisMaxNumberOfRecords, len(batch.Records))
}

func TestKinesis_ProcessBatchQueue_WillPostToKinesis(t *testing.T) {
	t.Skip("TODO")
}

func TestKinesis_ProcessBatchQueue_WillLogOnError(t *testing.T) {
	t.Skip("TODO")
}

func TestKinesis_ProcessBatchQueue_WillLogOnFailedRecords(t *testing.T) {
	t.Skip("TODO")
}

func Test_IsBatchReady(t *testing.T) {
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

	for _, testCase := range testCases {
		actuallyReady := isBatchReady(testCase.rl, testCase.ml, testCase.bc)
		assert.Equal(t, testCase.ready, actuallyReady)
	}
}

func TestGeneratePartitionKey(t *testing.T) {
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
