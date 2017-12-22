package dispatchers

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKinesisClient struct {
	mock.Mock
	kinesisiface.KinesisAPI
}

func (client *MockKinesisClient) PutRecords(input *kinesis.PutRecordsInput) (*kinesis.PutRecordsOutput, error) {
	args := client.Called(input)
	output, _ := args.Get(0).(*kinesis.PutRecordsOutput)
	err, _ := args.Get(1).(error)
	return output, err
}

func Test_KinesisService_New(t *testing.T) {
	svc := NewKinesisService("stream_name", "eu-west-1")
	assert.Equal(t, "stream_name", svc.streamName)
	assert.NotNil(t, svc.client)
}

func Test_KinesisService_CreateBatch(t *testing.T) {
	svc := &KinesisService{}
	assert.IsType(t, &KinesisBatch{}, svc.CreateBatch())
}

func Test_KinesisService_Send_ReturnsNil_OnSuccess(t *testing.T) {
	client := &MockKinesisClient{}
	svc := &KinesisService{
		streamName: "stream_name",
		client:     client,
	}
	batch := svc.CreateBatch()
	client.On("PutRecords", mock.AnythingOfType("*kinesis.PutRecordsInput")).
		Return(&kinesis.PutRecordsOutput{FailedRecordCount: aws.Int64(0)}, nil)
	err := svc.Send(batch)
	assert.Nil(t, err)
}

func Test_KinesisService_Send_ReturnsError_OnAWSError(t *testing.T) {
	client := &MockKinesisClient{}
	svc := &KinesisService{
		streamName: "stream_name",
		client:     client,
	}
	batch := svc.CreateBatch()
	client.On("PutRecords", mock.AnythingOfType("*kinesis.PutRecordsInput")).
		Return(nil, errors.New("AWS ERROR"))
	err := svc.Send(batch)
	assert.NotNil(t, err)
	assert.Equal(t, "AWS ERROR", err.Error())
}

func Test_KinesisService_Send_ReturnsError_OnFailedRecords(t *testing.T) {
	client := &MockKinesisClient{}
	svc := &KinesisService{
		streamName: "stream_name",
		client:     client,
	}
	batch := svc.CreateBatch()
	client.On("PutRecords", mock.AnythingOfType("*kinesis.PutRecordsInput")).
		Return(&kinesis.PutRecordsOutput{FailedRecordCount: aws.Int64(1)}, nil)
	err := svc.Send(batch)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed records")
}
