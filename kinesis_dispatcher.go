package main

type KinesisDispatcher struct {
	queue chan []byte
}

// Creates and returns a new dispatcher
func NewKinesisDispatcher(config *Config) *KinesisDispatcher {
	dispatcher := &KinesisDispatcher{
		queue: make(chan []byte, config.bufferSize),
	}
	return dispatcher
}

// inserts the message into the buffer for dispatching
// returns true for successful insertion, false if message was dropped
// will never block; if the queue is full, the message will be dropped
func (dispatcher *KinesisDispatcher) Put(message []byte) bool {
	select {
	case dispatcher.queue <- message:
		return true
	default:
		return false
	}
}

// fetches messages from the queue and dispatches to kinesis
func (dispatcher *KinesisDispatcher) Dispatch() {
	// TODO: think about batch processing
}
