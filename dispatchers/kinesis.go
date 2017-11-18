package dispatchers

// this is a stub struct for the moment

import (
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
)

type Kinesis struct {
	Session kinesisiface.KinesisAPI
}

func (client *Kinesis) Put(message []byte) bool {
	// TODO: implement functionality
	return true
}

func (client *Kinesis) Dispatch() {
}
