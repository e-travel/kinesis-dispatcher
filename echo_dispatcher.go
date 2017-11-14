package main

import "fmt"

// Used for local testing / debugging
type EchoDispatcher struct {
}

// prints the message to stdout
func (dispatcher *EchoDispatcher) Put(message []byte) {
	fmt.Println(string(message))
}

// Does nothing
func (dispatcher *EchoDispatcher) Dispatch() {
}
