package dispatchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageBuffer_Put_PlacesMessageInQueue(t *testing.T) {
	dispatcher := NewMessageBuffer(1, &MockDispatcher{})
	assert.True(t, dispatcher.Put([]byte("hello")))
	assert.Equal(t, []byte("hello"), <-dispatcher.queue)
}

func TestMessageBuffer_Put_IgnoresEmptyMessage(t *testing.T) {
	dispatcher := NewMessageBuffer(1, &MockDispatcher{})
	assert.False(t, dispatcher.Put([]byte("")))
	assert.Empty(t, dispatcher.queue)
}

func TestMessageBuffer_Put_DropsMessageWhenQueueIsFull(t *testing.T) {
	dispatcher := NewMessageBuffer(1, &MockDispatcher{})
	dispatcher.Put([]byte("hello"))
	assert.False(t, dispatcher.Put([]byte("goodbye")))
	assert.Equal(t, []byte("hello"), <-dispatcher.queue)
}

func TestMessageBuffer_Dispatch_WillForwardTheMessage(t *testing.T) {
	mockRecipient := &MockDispatcher{Messages: make(chan string)}
	dispatcher := NewMessageBuffer(1, mockRecipient)
	assert.True(t, dispatcher.Put([]byte("Hello There")))
	go dispatcher.Dispatch()
	assert.Equal(t, "Hello There", <-mockRecipient.Messages)
}
