package dispatchers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ======= MockService =======
type MockService struct {
	mock.Mock
}

func (svc *MockService) CreateBatch() Batch {
	args := svc.Called()
	b, _ := args.Get(0).(Batch)
	return b
}
func (svc *MockService) Send(batch Batch) error {
	args := svc.Called(batch)
	err, _ := args.Get(0).(error)
	return err
}

// ======= MockBatch =======
type MockBatch struct {
	mock.Mock
}

func (b *MockBatch) CanAdd(message []byte) bool {
	args := b.Called(message)
	ok, _ := args.Get(0).(bool)
	return ok
}
func (b *MockBatch) Len() int {
	args := b.Called()
	length, _ := args.Get(0).(int)
	return length
}

func TestMessageDispatcher_SetBatchFrequency(t *testing.T) {
	buffer := &MessageDispatcher{}
	buffer.SetBatchFrequency(100 * time.Second)
	assert.Equal(t, 100*time.Second, buffer.batchFrequency)
}

func TestMessageDispatcher_Put_PlacesMessageInQueue(t *testing.T) {
	buffer := NewMessageDispatcher(&MockService{}, 1)
	assert.True(t, buffer.Put([]byte("hello")))
	assert.Equal(t, []byte("hello"), <-buffer.messageQueue)
}

func TestMessageDispatcher_Put_IgnoresEmptyMessage(t *testing.T) {
	buffer := NewMessageDispatcher(&MockService{}, 1)
	assert.False(t, buffer.Put([]byte("")))
	assert.Empty(t, buffer.messageQueue)
}

func TestMessageDispatcher_Put_DropsMessageWhenQueueIsFull(t *testing.T) {
	buffer := NewMessageDispatcher(&MockService{}, 1)
	buffer.Put([]byte("hello"))
	assert.False(t, buffer.Put([]byte("goodbye")))
	assert.Equal(t, []byte("hello"), <-buffer.messageQueue)
}

func TestMessageDispatcher_Dispatch_WillProcessBothQueues(t *testing.T) {
	t.Skip("FIXME: This is broken")
	// setup
	service := &MockService{}
	batch := &MockBatch{}
	buffer := NewMessageDispatcher(service, 1)
	buffer.SetBatchFrequency(1 * time.Microsecond)
	// mock
	service.On("Send", batch).Return(nil)
	service.On("CreateBatch").Return(batch)
	batch.On("Len").Return(1)
	batch.On("CanAdd", mock.Anything).Return(true)

	// add message
	ok := buffer.Put([]byte("hello"))
	assert.True(t, ok)

	// fire
	buffer.Dispatch()
	// TODO: Use a proper wait condition instead of time.Sleep
	time.Sleep(2 * time.Second)
	batch.AssertExpectations(t)
	service.AssertExpectations(t)
}

func TestMessageDispatcher_processMessageQueue_WillAddMessageToBatch(t *testing.T) {
	t.Skip("FIXME: This is broken")
	// setup
	service := &MockService{}
	batch := &MockBatch{}
	buffer := NewMessageDispatcher(service, 2)
	message := []byte("Hello")
	// mock
	service.On("CreateBatch").Return(batch)
	batch.On("Len").Return(1)
	batch.On("CanAdd", message).Return(false)
	// add message
	ok := buffer.Put(message)
	assert.True(t, ok)
	// fire
	go buffer.processMessageQueue()
	timer := time.NewTimer(2 * time.Second)
	select {
	case <-buffer.batchQueue:
		assert.True(t, true)
	case <-timer.C:
		assert.Fail(t, "Timer expired")
	}
}

func TestMessageDispatcher_processMessageQueue_WillEnqueueBatch_OnTimer(t *testing.T) {
}

func TestMessageDispatcher_processMessageQueue_WillEnqueueBatch_OnReady(t *testing.T) {
}

func TestMessageDispatcher_processBatchQueue_WillDispatchBatch_UsingService(t *testing.T) {
}

func TestMessageDispatcher_processBatchQueue_WillIgnoreError(t *testing.T) {
}
