package main

import (
	"io"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockDispatcher struct {
	messages chan string
}

func (dispatcher *MockDispatcher) Put(message []byte) bool {
	dispatcher.messages <- string(message)
	return true
}

func (dispatcher *MockDispatcher) Dispatch() {
}

func DoRequest(config *Config, message string) error {
	conn, err := net.Dial(config.socketType, config.socketAddress)
	if err == nil {
		io.Copy(conn, strings.NewReader(message))
		conn.Close()
	}
	return err
}

func TestUnixSocketServerWillAcceptMessage(t *testing.T) {
	// setup
	config := &Config{
		socketType:    "unix",
		socketAddress: "/tmp/TestUnixSocketServerWillAcceptMessage.sock",
	}
	dispatcher := &MockDispatcher{
		messages: make(chan string),
	}
	running := make(chan bool)
	// start server
	go Serve(config, RequestHandler, dispatcher, running)
	defer os.Remove(config.socketAddress)
	<-running
	// fire
	err := DoRequest(config, "hello there")
	// check
	assert.Nil(t, err)
	assert.Equal(t, "hello there", <-dispatcher.messages)
}

func TestTCPSocketServerWillAcceptMessage(t *testing.T) {
	// setup
	config := &Config{
		socketType:    "tcp",
		socketAddress: "127.0.0.1:7778",
	}
	dispatcher := &MockDispatcher{
		messages: make(chan string),
	}
	running := make(chan bool)
	// start server
	go Serve(config, RequestHandler, dispatcher, running)
	<-running
	// fire
	err := DoRequest(config, "hello there")
	// check
	assert.Nil(t, err)
	assert.Equal(t, "hello there", <-dispatcher.messages)
}
