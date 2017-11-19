package servers

import "github.com/e-travel/message-dispatcher/dispatchers"

const (
	TCP      = "tcp"
	UNIX     = "unix"
	UNIXGRAM = "unixgram"
)

type Server interface {
	Serve(buffer dispatchers.Dispatcher, running chan<- bool)
}

func CreateServer(serverType string, serverAddress string) Server {
	var server Server
	switch serverType {
	case TCP, UNIX:
		server = &Stream{
			Type:    serverType,
			Address: serverAddress,
		}
	case UNIXGRAM:
		server = &UnixDatagram{
			Address: serverAddress,
		}
	}
	return server
}
