package servers

import (
	"net"
	"os"
	"testing"

	"github.com/e-travel/message-dispatcher/dispatchers"
	"github.com/stretchr/testify/assert"
)

func TestUnixDatagramServerWillAcceptMessages(t *testing.T) {
	// setup
	recipient := &dispatchers.MockDispatcher{Messages: make(chan string)}
	buffer := dispatchers.NewMessageBuffer(3, recipient)
	go buffer.Dispatch()
	running := make(chan bool)
	// start server
	server := &UnixDatagram{
		Address: "/tmp/TestUnixSocketServerWillAcceptMessage.sock",
	}
	go server.Serve(buffer, running)
	defer os.Remove(server.Address)
	<-running
	// open client connection
	conn, err := net.DialUnix(UNIXGRAM, nil,
		&net.UnixAddr{Name: server.Address, Net: UNIXGRAM})
	assert.Nil(t, err)
	defer conn.Close()
	// send messages
	go conn.Write([]byte("hello"))
	go conn.Write([]byte("there"))
	go conn.Write([]byte("bye"))
	// receive messages
	messages := make(map[string]bool)
	for i := 0; i < 3; i++ {
		messages[<-recipient.Messages] = true
	}
	// check
	assert.True(t, messages["hello"])
	assert.True(t, messages["there"])
	assert.True(t, messages["bye"])
}
