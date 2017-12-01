package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSocketTypeValidation(t *testing.T) {
	var testCases = []struct {
		value string
		valid bool
	}{
		{"tcp", false},
		{"unix", false},
		{"unixgram", true},
		{"udp", false},
		{"datagram", false},
		{"whatever", false},
		{"", false},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.valid, validateSocketType(testCase.value))
	}
}

func TestProducerTypeMapping(t *testing.T) {
	t.Skip("TODO: are values correctly mapped?")
}

func TestStreamNameValidation(t *testing.T) {
	var testCases = []struct {
		streamName     string
		dispatcherType string
		valid          bool
	}{
		{"", "kinesis", false},
		{"", "echo", true},
		{"a_name", "kinesis", true},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.valid,
			validateStreamName(testCase.streamName, testCase.dispatcherType))
	}
}
