package dispatchers

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockInfluxHttpClient struct {
	mock.Mock
}

func (c *MockInfluxHttpClient) WriteLines(lines *bytes.Buffer) error {
	args := c.Called(lines)
	err, _ := args.Get(0).(error)
	return err
}

func sendInfluxMessage(dispatcher *Influx, messageContent string) {
	dispatcher.Put([]byte(messageContent))
}

func fillInfluxMessageBuffer(dispatcher *Influx, messageContent string) {
	for i := 0; i < dispatcher.batchSize; i++ {
		sendInfluxMessage(dispatcher, messageContent)
	}
}

func TestInflux_Put_PlacesMessageToQueue(t *testing.T) {
	dispatcher := NewInflux(&MockInfluxHttpClient{})
	dispatcher.messageQueue = make(chan []byte, 1)
	assert.True(t, dispatcher.Put([]byte("hello")))
	assert.Equal(t, []byte("hello"), <-dispatcher.messageQueue)
}

func TestInflux_Put_DropsMessageWhenQueueIsFull(t *testing.T) {
	dispatcher := NewInflux(&MockInfluxHttpClient{})
	dispatcher.messageQueue = make(chan []byte, 1)
	dispatcher.Put([]byte("hello"))
	assert.False(t, dispatcher.Put([]byte("goodbye")))
	assert.Equal(t, []byte("hello"), <-dispatcher.messageQueue)
}

func TestInflux_Dispatch_WillProcessAllQueues(t *testing.T) {
	dispatcher := NewInflux(&MockInfluxHttpClient{})
	dispatcher.batchSize = 5
	fillInfluxMessageBuffer(dispatcher, "hello")
	sendInfluxMessage(dispatcher, "goodbye")
	dispatcher.Dispatch()
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

func TestInflux_processMessageQueue_WillAssembleBatch_WhenMaxSizeReached(t *testing.T) {
	dispatcher := NewInflux(&MockInfluxHttpClient{})
	dispatcher.batchSize = 5
	go dispatcher.processMessageQueue()
	// create a batch by filling the buffer
	fillInfluxMessageBuffer(dispatcher, "The same message all over again")
	// send one more message to trigger batch creation
	sendInfluxMessage(dispatcher, "This will stay in the queue")
	// get the batch
	<-dispatcher.batchQueue
}

func TestInflux_processMessageQueue_WillAssembleBatch_WhenTimerFires(t *testing.T) {
	dispatcher := NewInflux(&MockInfluxHttpClient{})
	dispatcher.maxBatchFrequency = time.Microsecond
	// send a single message to trigger batch creation
	sendInfluxMessage(dispatcher, "First message")
	go dispatcher.processMessageQueue()
	<-dispatcher.batchQueue
}

func TestInflux_processMessageQueue__WillDoNothing_WhenLinesEmpty(t *testing.T) {
	dispatcher := NewInflux(&MockInfluxHttpClient{})
	dispatcher.maxBatchFrequency = time.Microsecond
	// send a single empty message
	sendInfluxMessage(dispatcher, "")
	go dispatcher.processMessageQueue()
	assert.Empty(t, dispatcher.batchQueue)
}

func TestInflux_dispatchBatch_WillPutACopyToBatchQueue_AndResetBuffer(t *testing.T) {
	dispatcher := NewInflux(&MockInfluxHttpClient{})
	lines := bytes.NewBufferString("Hello")
	dispatcher.dispatchBatch(lines)
	// get the copy
	batch := <-dispatcher.batchQueue
	assert.Equal(t, 0, lines.Len())
	assert.NotEqual(t, lines, batch)
}

func TestInflux_processBatchQueue_WillSendBatchToInflux(t *testing.T) {
	// create mock client and inject to dispatcher
	client := &MockInfluxHttpClient{}
	dispatcher := NewInflux(client)
	dispatcher.client = client

	buf := bytes.NewBufferString("hello")

	// setup expectations
	client.On("WriteLines", buf).Once().Return(nil)

	// prepare
	dispatcher.batchQueue <- buf
	close(dispatcher.batchQueue)

	// fire
	dispatcher.processBatchQueue()

	client.AssertExpectations(t)
}
