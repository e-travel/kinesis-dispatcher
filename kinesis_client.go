package main

// this is a stub struct for the moment

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
)

type KinesisClient struct {
	svc kinesisiface.KinesisAPI
}

func NewKinesisClient() *KinesisClient {
	sess := session.Must(session.NewSession(&aws.Config{
		Retryer: client.DefaultRetryer{NumMaxRetries: 10},
		Region:  aws.String("eu-west-1"),
	}))
	return &KinesisClient{
		svc: kinesis.New(sess),
	}
}
