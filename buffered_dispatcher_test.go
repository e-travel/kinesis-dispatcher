package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockClient struct {
	messages chan string
}

func (client *MockClient) Put(message []byte) bool {
	client.messages <- string(message)
	return true
}

func TestPutPlacesMessageInQueue(t *testing.T) {
	config := &Config{
		bufferSize: 1,
	}
	dispatcher := NewBufferedDispatcher(config, &MockClient{})
	assert.True(t, dispatcher.Put([]byte("hello")))
	assert.Equal(t, []byte("hello"), <-dispatcher.queue)
}

func TestDropsMessageWhenQueueIsFull(t *testing.T) {
	config := &Config{
		bufferSize: 1,
	}
	dispatcher := NewBufferedDispatcher(config, &MockClient{})
	dispatcher.Put([]byte("hello"))
	assert.False(t, dispatcher.Put([]byte("goodbye")))
	assert.Equal(t, []byte("hello"), <-dispatcher.queue)
}

func TestMessageWillBeDispatchedToClient(t *testing.T) {
	config := &Config{
		bufferSize: 1,
	}
	mockClient := &MockClient{
		messages: make(chan string),
	}
	dispatcher := NewBufferedDispatcher(config, mockClient)
	assert.True(t, dispatcher.Put([]byte("Hello There")))
	go dispatcher.Dispatch()
	assert.Equal(t, "Hello There", <-mockClient.messages)
}
