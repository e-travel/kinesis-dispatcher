package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockClient struct{}

func (client *MockClient) Put(data [][]byte) bool {
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
