package main

import "fmt"

// Used for local testing / debugging
type EchoDispatcher struct {
}

// prints the message to stdout
func (dispatcher *EchoDispatcher) Put(message []byte) bool {
	fmt.Println(string(message))
	return true
}

// Does nothing
func (dispatcher *EchoDispatcher) Dispatch() {
}
