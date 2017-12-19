package dispatchers

// TODO: Remove this interface (we're covered by service now)
type Dispatcher interface {
	// accepts a message for dispatching
	Put([]byte) bool
	// the dispatching worker
	Dispatch()
}

// The following is a mock type for use in tests
type MockDispatcher struct {
	Messages chan string
}

func (dispatcher *MockDispatcher) Put(message []byte) bool {
	dispatcher.Messages <- string(message)
	return true
}
func (dispatcher *MockDispatcher) Dispatch() {
	for range dispatcher.Messages {
	}
}
