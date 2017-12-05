package main

import (
	"testing"
	"time"

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
		{"udp", true},
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

func TestInfluxValidation(t *testing.T) {
	var testCases = []struct {
		dispatcherType string
		host           string
		database       string
		valid          bool
	}{
		{"influx", "", "", false},
		{"influx", "host", "", false},
		{"influx", "", "database", false},
		{"influx", "host", "database", true},
		{"other", "", "", true},
		{"", "", "", true},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.valid,
			validateInflux(testCase.dispatcherType, testCase.host, testCase.database))
	}
}

func TestBatchFrequencyValidation(t *testing.T) {
	var testCases = []struct {
		freq  time.Duration
		valid bool
	}{
		{-1 * time.Second, false},
		{0 * time.Second, false},
		{1 * time.Second, true},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.valid, validateBatchFrequency(testCase.freq))
	}
}
