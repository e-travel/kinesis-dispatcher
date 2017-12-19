package dispatchers

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
)

type KinesisService struct {
	client     kinesisiface.KinesisAPI
	streamName string
}

func NewKinesisService(streamName string, awsRegion string) *KinesisService {
	// create session
	sess := session.Must(session.NewSession(&aws.Config{
		Retryer: client.DefaultRetryer{NumMaxRetries: 10},
		Region:  aws.String(awsRegion),
	}))
	return &KinesisService{
		client:     kinesis.New(sess),
		streamName: streamName,
	}
}

func (svc *KinesisService) CreateBatch() Batch {
	return NewKinesisBatch(svc.streamName)
}

func (svc *KinesisService) Send(batch Batch) error {
	b := batch.(*KinesisBatch)
	output, err := svc.client.PutRecords(b.records)
	if err != nil {
		return err
	}
	if *output.FailedRecordCount > 0 {
		return errors.New(fmt.Sprintf("AWS Kinesis: failed records %d/%d",
			*output.FailedRecordCount, len(b.records.Records)))
	}
	return nil
}
