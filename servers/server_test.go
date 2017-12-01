package servers

import (
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/e-travel/message-dispatcher/dispatchers"
	"github.com/stretchr/testify/assert"
)

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
