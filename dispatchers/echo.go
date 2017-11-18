package dispatchers

import "fmt"

// Used for local testing / debugging
type Echo struct {
}

// prints the message to stdout
func (dispatcher *Echo) Put(message []byte) bool {
	fmt.Println(string(message))
	return true
}

// Does nothing
func (dispatcher *Echo) Dispatch() {
}
