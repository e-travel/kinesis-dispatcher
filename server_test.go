package main

import (
	"io"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/e-travel/message-dispatcher/dispatchers"
	"github.com/stretchr/testify/assert"
)

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
	recipient := &dispatchers.MockDispatcher{Messages: make(chan string)}
	buffer := dispatchers.NewMessageBuffer(config.bufferSize, recipient)
	go buffer.Dispatch()
	running := make(chan bool)
	// start server
	go Serve(config, RequestHandler, buffer, running)
	defer os.Remove(config.socketAddress)
	<-running
	// fire
	err := DoRequest(config, "hello there")
	// check
	assert.Nil(t, err)
	assert.Equal(t, "hello there", <-recipient.Messages)
}

func TestTCPSocketServerWillAcceptMessage(t *testing.T) {
	// setup
	config := &Config{
		socketType:    "tcp",
		socketAddress: "127.0.0.1:7778",
	}
	recipient := &dispatchers.MockDispatcher{Messages: make(chan string)}
	buffer := dispatchers.NewMessageBuffer(config.bufferSize, recipient)
	go buffer.Dispatch()
	running := make(chan bool)
	// start server
	go Serve(config, RequestHandler, buffer, running)
	<-running
	// fire
	err := DoRequest(config, "hello there")
	// check
	assert.Nil(t, err)
	assert.Equal(t, "hello there", <-recipient.Messages)
}
