package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutPlacesMessageInQueue(t *testing.T) {
	config := &Config{
		bufferSize: 1,
	}
	dispatcher := NewBufferedDispatcher(config, &KinesisClient{})
	assert.True(t, dispatcher.Put([]byte("hello")))
	assert.Equal(t, []byte("hello"), <-dispatcher.queue)
}

func TestDropsMessageWhenQueueIsFull(t *testing.T) {
	config := &Config{
		bufferSize: 1,
	}
	dispatcher := NewBufferedDispatcher(config, &KinesisClient{})
	dispatcher.Put([]byte("hello"))
	assert.False(t, dispatcher.Put([]byte("goodbye")))
	assert.Equal(t, []byte("hello"), <-dispatcher.queue)
}
