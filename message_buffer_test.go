package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutPlacesMessageInQueue(t *testing.T) {
	config := &Config{
		bufferSize: 1,
	}
	dispatcher := NewMessageBuffer(config, &MockDispatcher{})
	assert.True(t, dispatcher.Put([]byte("hello")))
	assert.Equal(t, []byte("hello"), <-dispatcher.queue)
}

func TestPutDropsMessageWhenQueueIsFull(t *testing.T) {
	config := &Config{
		bufferSize: 1,
	}
	dispatcher := NewMessageBuffer(config, &MockDispatcher{})
	dispatcher.Put([]byte("hello"))
	assert.False(t, dispatcher.Put([]byte("goodbye")))
	assert.Equal(t, []byte("hello"), <-dispatcher.queue)
}

func TestDispatchWillForwardTheMessage(t *testing.T) {
	config := &Config{
		bufferSize: 1,
	}
	mockRecipient := NewMockDispatcher(1)
	dispatcher := NewMessageBuffer(config, mockRecipient)
	assert.True(t, dispatcher.Put([]byte("Hello There")))
	go dispatcher.Dispatch()
	assert.Equal(t, "Hello There", <-mockRecipient.messages)
}
