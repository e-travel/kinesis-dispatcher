package servers

import (
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/e-travel/message-dispatcher/dispatchers"
)

type Stream struct {
	Type    string
	Address string
}

// TODO: what if this function panics?
func HandleStream(conn net.Conn, buffer dispatchers.Dispatcher) {
	defer conn.Close()
	b, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Print("Error reading from connection")
		return
	}
	// TODO: do some logging here if Put returns false
	buffer.Put(b)
}

func (server *Stream) Serve(buffer dispatchers.Dispatcher, running chan<- bool) {
	// remove any existing socket file
	if server.Type == "unix" {
		os.Remove(server.Address)
	}
	listener, err := net.Listen(server.Type, server.Address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	running <- true

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go HandleStream(conn, buffer)
	}
}
