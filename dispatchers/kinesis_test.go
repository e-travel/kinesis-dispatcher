package dispatchers

import (
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
	<-dispatcher.batchQueue
	assert.Empty(t, dispatcher.messageQueue)
	assert.Empty(t, dispatcher.batchQueue)
}

func TestKinesis_processMessageQueue_WillPutInBatchQueue_WhenReady(t *testing.T) {
	dispatcher := NewKinesis("stream_name", "region", 10*time.Second)
	batch := NewKinesisBatch("stream_name")
	go dispatcher.processMessageQueue(batch)
	// create a batch by filling the buffer
	fillMessageBuffer(dispatcher, "The same message all over again")
	// send one more message to trigger batch creation
	sendMessage(dispatcher, "This will stay in the queue")
	// get the batch
	batch = <-dispatcher.batchQueue
	assert.Equal(t, KinesisMaxNumberOfRecords, batch.Len())
}

func TestKinesis_processMessageQueue_WillDropMessage_WhenTooLarge(t *testing.T) {
	dispatcher := NewKinesis("stream_name", "region", 10*time.Second)
	batch := NewKinesisBatch("stream_name")
	go dispatcher.processMessageQueue(batch)
	// create a batch by filling the buffer
	fillMessageBuffer(dispatcher, "The same message all over again")
	// send one more message to trigger batch creation
	sendMessage(dispatcher, "This will stay in the queue")
	// get the batch
	batch = <-dispatcher.batchQueue
	assert.Equal(t, KinesisMaxNumberOfRecords, batch.Len())
}

func TestKinesis_processMessageQueue_WillPutInBatchQueue_WhenTimerFires(t *testing.T) {
	//t.Skip("Implement timer test")
	dispatcher := NewKinesis("stream_name", "region", time.Microsecond)
	batch := NewKinesisBatch("stream_name")
	sendMessage(dispatcher, "message")
	go dispatcher.processMessageQueue(batch)
	// get the batch
	batch = <-dispatcher.batchQueue
	assert.Equal(t, 1, batch.Len())
}

func TestKinesis_processBatchQueue_WillSendBatchToKinesis(t *testing.T) {
	t.Skip("TODO")
}

func TestKinesis_send_WillPostBatchToKinesis(t *testing.T) {
	t.Skip("TODO")
}

func TestKinesis_send_WillReturnErrorOnKinesisError(t *testing.T) {
	t.Skip("TODO")
}

func TestKinesis_send_WillReturnErrorOnFailedRecords(t *testing.T) {
	t.Skip("TODO")
}
