package servers

import (
	"log"
	"net"

	"github.com/e-travel/message-dispatcher/dispatchers"
)

type Udp struct {
	Address string
}

func (server *Udp) Serve(buffer dispatchers.Dispatcher, running chan<- bool) {
	udpAddr, err := net.ResolveUDPAddr("udp", server.Address)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.ListenUDP(UDP, udpAddr)
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
