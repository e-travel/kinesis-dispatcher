package dispatchers

type MessageBuffer struct {
	queue     chan []byte
	recipient Dispatcher
}

// Creates and returns a new dispatcher
func NewMessageBuffer(bufferSize int, recipient Dispatcher) *MessageBuffer {
	return &MessageBuffer{
		queue:     make(chan []byte, bufferSize),
		recipient: recipient,
	}
}

// inserts the message into the buffer for dispatching
// returns true for successful insertion, false if message was dropped
// will never block; if the queue is full, the message will be dropped
func (dispatcher *MessageBuffer) Put(message []byte) bool {
	select {
	case dispatcher.queue <- message:
		return true
	default:
		return false
	}
}

// fetches messages from the queue and dispatches to recipient
func (dispatcher *MessageBuffer) Dispatch() {
	for message := range dispatcher.queue {
		dispatcher.recipient.Put(message)
	}
}
