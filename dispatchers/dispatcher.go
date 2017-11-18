package dispatchers

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
func (_ *MockDispatcher) Dispatch() {}
