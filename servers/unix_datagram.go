package servers

import (
	"net"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/e-travel/message-dispatcher/dispatchers"
)

type UnixDatagram struct {
	Address string
}

func (server *UnixDatagram) Serve(buffer dispatchers.Dispatcher, running chan<- bool) {
	os.Remove(server.Address)
	conn, err := net.ListenUnixgram(
		UNIXGRAM,
		&net.UnixAddr{
			Name: server.Address,
			Net:  UNIXGRAM,
		})
	if err != nil {
		log.Fatal(err)
	}

	running <- true
	for {
		var b [8192]byte
		n, err := conn.Read(b[:])
		if err != nil {
			log.Fatal(err)
		}
		buffer.Put(b[:n])
	}
}
