package dispatchers

// Used for local testing / debugging
type Echo struct {
}

// prints the message to stdout
func (dispatcher *Echo) Put(message []byte) bool {
	return true
}

// Does nothing
func (dispatcher *Echo) Dispatch() {
}
