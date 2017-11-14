package main

type Dispatcher interface {
	// accepts a message for dispatching
	Put([]byte)
	// the dispatching worker
	Dispatch()
}
