package servers

import (
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/e-travel/message-dispatcher/dispatchers"
	"github.com/stretchr/testify/assert"
)

func TestCreateServerReturnsStreamWhenTCP(t *testing.T) {
	server := CreateServer(TCP, ":8888").(*Stream)
	assert.Equal(t, TCP, server.Type)
	assert.Equal(t, ":8888", server.Address)
}

func TestCreateServerReturnsStreamWhenUNIX(t *testing.T) {
	server := CreateServer(UNIX, "/tmp/sock").(*Stream)
	assert.Equal(t, UNIX, server.Type)
	assert.Equal(t, "/tmp/sock", server.Address)
}

func TestCreateServerReturnsUnixDatagramWhenUNIXGRAM(t *testing.T) {
	server := CreateServer(UNIXGRAM, "/tmp/sock").(*UnixDatagram)
	assert.Equal(t, "/tmp/sock", server.Address)
}

// ==========
// Benchmarks
// ==========

type NullDispatcher struct{}

func (_ *NullDispatcher) Put(_ []byte) bool {
	return true
}
func (_ *NullDispatcher) Dispatch() {
}

type MockBlockingMessageBuffer struct {
	queue     chan []byte
	recipient dispatchers.Dispatcher
	dispatchers.MessageBuffer
}

func (buf *MockBlockingMessageBuffer) Put(message []byte) bool {
	buf.queue <- message
	return true
}

func BenchmarkSendMessageToUnixDatagramServer(b *testing.B) {
	unixDatagramSocket := "/tmp/BenchmarkUnixDatagramServer.sock"
	defer os.Remove(unixDatagramSocket)
	server := &UnixDatagram{Address: unixDatagramSocket}
	dispatcher := &NullDispatcher{}
	buffer := &MockBlockingMessageBuffer{
		queue:     make(chan []byte, 1024),
		recipient: dispatcher,
	}
	go buffer.Dispatch()
	running := make(chan bool)
	go server.Serve(buffer, running)
	<-running
	conn, err := net.Dial("unixgram", unixDatagramSocket)
	if err != nil {
		b.Fatal(err.Error())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Fprintf(conn, "Payload #%d", i)
	}
}

// ========
// Examples
// ========

func ExampleSendMessageToUnixStreamServer() {
	socket := "/tmp/ExampleSendMessageToUnixStreamServer.sock"
	defer os.Remove(socket)
	server := &Stream{Type: "unix", Address: socket}
	dispatcher := &NullDispatcher{}
	buffer := &MockBlockingMessageBuffer{
		queue:     make(chan []byte, 1024),
		recipient: dispatcher,
	}
	go buffer.Dispatch()
	running := make(chan bool)
	go server.Serve(buffer, running)
	<-running
	conn, err := net.Dial(server.Type, server.Address)
	if err != nil {
		fmt.Printf("Could not connect to %s", server.Address)
		return
	}
	fmt.Fprintf(conn, "Payload\n")
	conn.Close()
	// Output:
}

func ExampleSendMessageToTCPServer() {
	server := &Stream{Type: "tcp", Address: "127.0.0.1:11234"}
	dispatcher := &NullDispatcher{}
	buffer := &MockBlockingMessageBuffer{
		queue:     make(chan []byte, 1024),
		recipient: dispatcher,
	}
	go buffer.Dispatch()
	running := make(chan bool)
	go server.Serve(buffer, running)
	<-running
	conn, err := net.Dial(server.Type, server.Address)
	if err != nil {
		fmt.Printf("Could not connect to %s", server.Address)
		return
	}
	conn.Close()
	// Output:
}
