package main

import (
	"flag"
	"os"

	"github.com/e-travel/message-dispatcher/servers"
)

type Config struct {
	socketType    string
	socketAddress string
	streamName    string
	bufferSize    int
	echoMode      bool
}

var ValidSocketTypes = map[string]bool{
	servers.TCP:      true,
	servers.UNIX:     true,
	servers.UNIXGRAM: true,
}

func ParseFromCommandLine(config *Config) {
	flag.StringVar(&config.socketType, "type", servers.UNIXGRAM, "The socket's type (tcp, unix, unixgram)")
	flag.StringVar(&config.socketAddress, "address", ":8888", "The socket's address (port or file)")
	flag.StringVar(&config.streamName, "streamName", "", "The name of the kinesis stream")
	flag.IntVar(&config.bufferSize, "size", 1024, "The size of the buffer")
	flag.BoolVar(&config.echoMode, "echo", false, "Activate echo mode (no kinesis requests)")
	helpRequested := flag.Bool("help", false, "Print usage help and exit")

	flag.Parse()

	if *helpRequested {
		flag.Usage()
		os.Exit(0)
	}
}

func (config *Config) Validate() bool {
	return validateSocketType(config.socketType) &&
		validateStreamName(config.streamName)
}

func validateSocketType(socketType string) bool {
	return ValidSocketTypes[socketType]
}

func validateStreamName(streamName string) bool {
	return len(streamName) > 0
}
