package servers

import (
	"net"
	"testing"

	"github.com/e-travel/message-dispatcher/dispatchers"
	"github.com/stretchr/testify/assert"
)

func TestUdpServerWillAcceptMessages(t *testing.T) {
	// setup server
	recipient := &dispatchers.MockDispatcher{Messages: make(chan string)}
	buffer := dispatchers.NewMessageBuffer(3, recipient)
	go buffer.Dispatch()
	running := make(chan bool)
	// start server
	server := &Udp{
		Address: ":9999",
	}
	go server.Serve(buffer, running)
	<-running

	// setup client
	conn, err := net.Dial("udp", "127.0.0.1:9999")
	assert.Nil(t, err)
	defer conn.Close()
	conn.Write([]byte("hello"))
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
	assert.Equal(t, 3, len(messages))
}
