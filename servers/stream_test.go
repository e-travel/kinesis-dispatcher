package servers

import (
	"io"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/e-travel/message-dispatcher/dispatchers"
	"github.com/stretchr/testify/assert"
)

func DoRequest(server *Stream, message string) error {
	conn, err := net.Dial(server.Type, server.Address)
	if err == nil {
		io.Copy(conn, strings.NewReader(message))
		conn.Close()
	}
	return err
}

func TestUnixSocketServerWillAcceptMessage(t *testing.T) {
	// setup
	recipient := &dispatchers.MockDispatcher{Messages: make(chan string)}
	buffer := dispatchers.NewMessageBuffer(1, recipient)
	go buffer.Dispatch()
	running := make(chan bool)
	// start server
	server := &Stream{
		Type:    "unix",
		Address: "/tmp/TestUnixSocketServerWillAcceptMessage.sock",
	}
	go server.Serve(buffer, running)
	defer os.Remove(server.Address)
	<-running
	// fire
	err := DoRequest(server, "hello there")
	// check
	assert.Nil(t, err)
	assert.Equal(t, "hello there", <-recipient.Messages)
}

func TestTCPSocketServerWillAcceptMessage(t *testing.T) {
	// setup
	recipient := &dispatchers.MockDispatcher{Messages: make(chan string)}
	buffer := dispatchers.NewMessageBuffer(1, recipient)
	go buffer.Dispatch()
	running := make(chan bool)
	// start server
	server := &Stream{
		Type:    "tcp",
		Address: "127.0.0.1:7778",
	}
	go server.Serve(buffer, running)
	<-running
	// fire
	err := DoRequest(server, "hello there")
	// check
	assert.Nil(t, err)
	assert.Equal(t, "hello there", <-recipient.Messages)
}
