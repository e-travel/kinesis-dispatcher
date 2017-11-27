package main

import (
	"testing"

	"github.com/e-travel/message-dispatcher/dispatchers"
	"github.com/stretchr/testify/assert"
)

func TestCreateDispatcher(t *testing.T) {
	var testCases = []struct {
		name   string
		errMsg string
		obj    dispatchers.Dispatcher
	}{
		{"echo", "", &dispatchers.Echo{}},
		{"kinesis", "", &dispatchers.Kinesis{}},
		{"foo", "Invalid dispatcher type: foo", nil},
	}
	for _, testCase := range testCases {
		dispatcher, err := createDispatcher(testCase.name)
		assert.Equal(t, testCase.errMsg == "", err == nil)
		if err != nil {
			assert.Equal(t, testCase.errMsg, err.Error())
		}
		assert.IsType(t, testCase.obj, dispatcher)
	}
}
