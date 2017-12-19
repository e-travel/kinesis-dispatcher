package servers

import (
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/e-travel/message-dispatcher/dispatchers"
	"github.com/stretchr/testify/assert"
)

func TestCreateServer_ReturnsUnixgram_OnRequest(t *testing.T) {
	server, ok := CreateServer(UNIXGRAM, "/tmp/sock").(*UnixDatagram)
	assert.NotNil(t, ok)
	assert.Equal(t, "/tmp/sock", server.Address)
}

func TestCreateServer_ReturnsUDP_OnRequest(t *testing.T) {
	server, ok := CreateServer(UDP, ":8888").(*Udp)
	assert.NotNil(t, ok)
	assert.Equal(t, ":8888", server.Address)
}

// ==========
// Benchmarks
// ==========

func BenchmarkSendMessageToUnixDatagramServer(b *testing.B) {
	unixDatagramSocket := "/tmp/BenchmarkUnixDatagramServer.sock"
	defer os.Remove(unixDatagramSocket)
	server := &UnixDatagram{Address: unixDatagramSocket}
	dispatcher := &dispatchers.MockDispatcher{Messages: make(chan string, 1024)}
	go dispatcher.Dispatch()
	running := make(chan bool)
	go server.Serve(dispatcher, running)
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
