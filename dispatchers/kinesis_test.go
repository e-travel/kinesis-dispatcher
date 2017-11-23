package dispatchers

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/stretchr/testify/assert"
)

func fillMessageBuffer(dispatcher *Kinesis, messageContent string) {
	for i := 0; i < KinesisMaxNumberOfRecords; i++ {
		dispatcher.Put([]byte(messageContent))
	}
}

func TestKinesisConstants(t *testing.T) {
	assert.Equal(t, 500, KinesisMaxNumberOfRecords)
	assert.Equal(t, 5*MEGABYTE, KinesisMaxSizeInBytes)
	assert.Equal(t, 1000, KinesisBufferSize)
}

func TestKinesisPutPlacesMessageToQueue(t *testing.T) {
	dispatcher := NewKinesis("stream_name")
	dispatcher.messageQueue = make(chan []byte, 1)
	assert.True(t, dispatcher.Put([]byte("hello")))
	assert.Equal(t, []byte("hello"), <-dispatcher.messageQueue)
}

func TestKinesisPutMessageWhenQueueIsFull(t *testing.T) {
	dispatcher := NewKinesis("stream_name")
	dispatcher.messageQueue = make(chan []byte, 1)
	dispatcher.Put([]byte("hello"))
	assert.False(t, dispatcher.Put([]byte("goodbye")))
	assert.Equal(t, []byte("hello"), <-dispatcher.messageQueue)
}

func TestKinesisDispatchWillProcessAllQueues(t *testing.T) {
	dispatcher := NewKinesis("stream_name")
	fillMessageBuffer(dispatcher, "hello")
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

func TestKinesisProcessMessageQueueWillAssembleBatchAndPutInBatchQueue(t *testing.T) {
	dispatcher := NewKinesis("stream_name")
	go dispatcher.processMessageQueue()
	// create a batch by filling the buffer
	fillMessageBuffer(dispatcher, "The same message all over again")
	// get the batch
	batch := <-dispatcher.batchQueue
	assert.Equal(t, KinesisMaxNumberOfRecords, len(batch.Records))
	assert.Empty(t, dispatcher.messageQueue)
}

func TestKinesisProcessMessageQueueWillNotBlockWhenEnqueingBatch(t *testing.T) {
	dispatcher := NewKinesis("stream_name")
	// downsize the batch queue
	dispatcher.batchQueue = make(chan *kinesis.PutRecordsInput, 1)
	// create two batches
	fillMessageBuffer(dispatcher, "The same message all over again")
	fillMessageBuffer(dispatcher, "Another message all over again")
	close(dispatcher.messageQueue)
	// TODO: the success of this test depends on the following statement not blocking
	// this is not a good practice since the whole test suite may block
	// we need to add Context in Dispatch (and the process*Queue methods)
	dispatcher.processMessageQueue()
	<-dispatcher.batchQueue
	// no batch is left to process
	assert.Empty(t, dispatcher.batchQueue)
}

func TestKinesisProcessBatchQueueWillPostToKinesis(t *testing.T) {
}

func TestKinesisProcessBatchQueueWillLogOnError(t *testing.T) {
}

func TestKinesisProcessBatchQueueWillLogOnFailedRecords(t *testing.T) {
}

func TestIsBatchReady(t *testing.T) {
	var testCases = []struct {
		rs    int // records size
		bc    int // byte count
		ml    int // message length
		ready bool
	}{
		{0, 0, 10, false},
		{400, 1 * MEGABYTE, 100, false},
		{500, 1 * MEGABYTE, 100, true},
		{400, 5*MEGABYTE - 1, 0, false},
		{400, 5*MEGABYTE - 1, 1, true},
	}

	for _, testCase := range testCases {
		actuallyReady := isBatchReady(testCase.ml, testCase.rs, testCase.bc)
		assert.Equal(t, testCase.ready, actuallyReady)
	}
}
