package main

type Dispatcher interface {
	// accepts a message for dispatching
	Put([]byte) bool
	// the dispatching worker
	Dispatch()
}
