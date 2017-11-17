package main

type BufferedDispatcher struct {
	queue  chan []byte
	client Client
}

// Creates and returns a new dispatcher
func NewBufferedDispatcher(config *Config, client Client) *BufferedDispatcher {
	dispatcher := &BufferedDispatcher{
		queue:  make(chan []byte, config.bufferSize),
		client: client,
	}
	return dispatcher
}

// inserts the message into the buffer for dispatching
// returns true for successful insertion, false if message was dropped
// will never block; if the queue is full, the message will be dropped
func (dispatcher *BufferedDispatcher) Put(message []byte) bool {
	select {
	case dispatcher.queue <- message:
		return true
	default:
		return false
	}
}

// fetches messages from the queue and dispatches to kinesis
func (dispatcher *BufferedDispatcher) Dispatch() {
	// TODO: think about batch processing
}
