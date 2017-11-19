package servers

import (
	"testing"

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
