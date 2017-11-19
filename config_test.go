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

	config := &Config{}

	for _, testCase := range testCases {
		config.socketType = testCase.value
		assert.Equal(t, testCase.valid, config.Validate())
	}
}
