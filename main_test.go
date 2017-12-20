package main

import (
	"testing"

	"github.com/e-travel/message-dispatcher/dispatchers"
	"github.com/stretchr/testify/assert"
)

func TestCreateDispatcher(t *testing.T) {
	var testCases = []struct {
		config Config
		errMsg string
		svc    dispatchers.Service
	}{
		{Config{dispatcherType: "echo"}, "", &dispatchers.EchoService{}},
		{Config{dispatcherType: "kinesis"}, "", &dispatchers.KinesisService{}},
		{Config{dispatcherType: "influx"}, "", &dispatchers.InfluxService{}},
		{Config{dispatcherType: "foo"}, "Invalid backend type: foo", nil},
	}
	for _, testCase := range testCases {
		dispatcher, err := createDispatcher(&testCase.config)
		assert.Equal(t, testCase.errMsg == "", err == nil)
		if err != nil {
			assert.Equal(t, testCase.errMsg, err.Error())
		}
		assert.IsType(t, testCase.svc, dispatcher.Service)
	}
}
