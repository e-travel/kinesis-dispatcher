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
	b, _ := args.Get(0).(*MockBatch)
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

func (b *MockBatch) Add(message []byte) error {
	args := b.Called(message)
	err, _ := args.Get(0).(error)
	return err
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

// ======== TESTS ========
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

func TestMessageDispatcher_processMessageQueue_WillAddMessageToBatch(t *testing.T) {
	// setup
	service := &MockService{}
	batch := &MockBatch{}

	buffer := NewMessageDispatcher(service, 2)
	message := []byte("Hello")
	// mock
	service.On("CreateBatch").Return(batch)
	batch.On("Len").Return(1)
	batch.On("CanAdd", message).Return(false)
	batch.On("Add", mock.Anything).Return(nil)
	// add message
	ok := buffer.Put(message)
	assert.True(t, ok)
	close(buffer.messageQueue)
	// fire
	buffer.processMessageQueue()
	<-buffer.batchQueue
	batch.AssertExpectations(t)
	service.AssertExpectations(t)
}

func TestMessageDispatcher_processMessageQueue_WillEnqueueBatch_OnTimer(t *testing.T) {
	// setup
	service := &MockService{}
	batch := &MockBatch{}

	buffer := NewMessageDispatcher(service, 2)
	buffer.SetBatchFrequency(time.Millisecond)
	message := []byte("Hello")
	// mock
	service.On("CreateBatch").Return(batch)
	batch.On("Len").Return(1)
	batch.On("CanAdd", message).Return(false)
	batch.On("Add", mock.Anything).Return(nil)
	// add message
	ok := buffer.Put(message)
	assert.True(t, ok)
	close(buffer.messageQueue)
	// fire
	buffer.processMessageQueue()
	<-buffer.batchQueue
	batch.AssertExpectations(t)
	service.AssertExpectations(t)
}

func TestMessageDispatcher_processBatchQueue_WillDispatchBatch_UsingService(t *testing.T) {
	// setup
	service := &MockService{}
	buffer := NewMessageDispatcher(service, 2)
	batch := &MockBatch{}
	buffer.batchQueue <- batch
	close(buffer.batchQueue)
	// expectation
	service.On("Send", batch).Once().Return(nil)
	// fire
	buffer.processBatchQueue()
	service.AssertExpectations(t)
}
