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
		{"tcp", true},
		{"unix", true},
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

func TestStreamNameValidation(t *testing.T) {
	var testCases = []struct {
		value string
		valid bool
	}{
		{"", false},
		{"a_name", true},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.valid, validateStreamName(testCase.value))
	}
}
