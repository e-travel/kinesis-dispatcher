package main

// this is a stub struct for the moment

import (
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
)

type KinesisClient struct {
	session kinesisiface.KinesisAPI
}

func (client *KinesisClient) Put(message []byte) bool {
	// TODO: implement functionality
	return true
}
