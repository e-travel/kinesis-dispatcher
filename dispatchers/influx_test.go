package dispatchers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func sendMessage(dispatcher *Influx, messageContent string) {
	dispatcher.Put([]byte(messageContent))
}

func fillMessageBuffer(dispatcher *Influx, messageContent string) {
	for i := 0; i < InfluxMaxNumberOfRecords; i++ {
		sendMessage(dispatcher, messageContent)
	}
}

func TestInfluxConstants(t *testing.T) {
	assert.Equal(t, 500, InfluxMaxNumberOfRecords)
	assert.Equal(t, 5*MEGABYTE, InfluxMaxSizeInBytes)
	assert.Equal(t, 1000, InfluxBufferSize)
}

func TestInfluxPutPlacesMessageToQueue(t *testing.T) {
	dispatcher := NewInflux("stream_name", "region")
	dispatcher.messageQueue = make(chan []byte, 1)
	assert.True(t, dispatcher.Put([]byte("hello")))
	assert.Equal(t, []byte("hello"), <-dispatcher.messageQueue)
}

func TestInfluxPutMessageWhenQueueIsFull(t *testing.T) {
	dispatcher := NewInflux("stream_name", "region")
	dispatcher.messageQueue = make(chan []byte, 1)
	dispatcher.Put([]byte("hello"))
	assert.False(t, dispatcher.Put([]byte("goodbye")))
	assert.Equal(t, []byte("hello"), <-dispatcher.messageQueue)
}

func TestInfluxDispatchWillProcessAllQueues(t *testing.T) {
	dispatcher := NewInflux("stream_name", "region")
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

func TestInfluxProcessMessageQueueWillAssembleBatchAndPutInBatchQueue(t *testing.T) {
	//t.Skip("This blocks; needs fixing")
	dispatcher := NewInflux("stream_name", "region")
	go dispatcher.processMessageQueue()
	// create a batch by filling the buffer
	fillMessageBuffer(dispatcher, "The same message all over again")
	// send one more message to trigger batch creation
	sendMessage(dispatcher, "This will stay in the queue")
	// get the batch
	batch := <-dispatcher.batchQueue
	assert.Equal(t, InfluxMaxNumberOfRecords, len(batch.Records))
}

func TestInfluxProcessBatchQueueWillPostToInflux(t *testing.T) {
	t.Skip("TODO")
}

func TestInfluxProcessBatchQueueWillLogOnError(t *testing.T) {
	t.Skip("TODO")
}

func TestInfluxProcessBatchQueueWillLogOnFailedRecords(t *testing.T) {
	t.Skip("TODO")
}
