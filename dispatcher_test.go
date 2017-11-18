package main

// defines a mock type for use in other tests

type MockDispatcher struct {
	messages chan string
}

func (dispatcher *MockDispatcher) Put(message []byte) bool {
	dispatcher.messages <- string(message)
	return true
}

func (_ *MockDispatcher) Dispatch() {}

func NewMockDispatcher(size int) *MockDispatcher {
	return &MockDispatcher{
		messages: make(chan string, size),
	}
}
