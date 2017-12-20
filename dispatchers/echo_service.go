package dispatchers

import "fmt"

// Used for local testing / debugging
type EchoService struct{}

// prints the message to stdout
func (svc *EchoService) CreateBatch() Batch {
	return &EchoBatch{}
}

func (svc *EchoService) Send(batch Batch) error {
	b := batch.(*EchoBatch)
	for _, msg := range b.Messages {
		fmt.Printf("msg=%s\n", string(msg))
	}
	return nil
}
