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
func (dispatcher *KinesisDispatcher) Put(message []byte) {
	select {
	case dispatcher.queue <- message:
	default:
		// do not block if channel is full
		// this means we may drop messages
	}
}

// fetches messages from the queue and dispatches to kinesis
func (dispatcher *KinesisDispatcher) Dispatch() {
	// TODO: think about batch processing
}
