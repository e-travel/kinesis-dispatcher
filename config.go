package main

import (
	"flag"
	"os"

	"github.com/e-travel/message-dispatcher/servers"
)

type Config struct {
	socketType     string
	socketAddress  string
	streamName     string
	awsRegion      string
	bufferSize     int
	dispatcherType string
}

var ValidSocketTypes = map[string]bool{
	servers.TCP:      true,
	servers.UNIX:     true,
	servers.UNIXGRAM: true,
}

func ParseFromCommandLine(config *Config) {
	flag.StringVar(&config.socketType, "type", servers.UNIXGRAM, "The socket's type (tcp, unix, unixgram)")
	flag.StringVar(&config.socketAddress, "address", ":8888", "The socket's address (port or file)")
	flag.StringVar(&config.streamName, "stream-name", "", "The name of the kinesis stream")
	flag.StringVar(&config.awsRegion, "aws-region", "eu-west-1", "The kinesis stream's AWS region")
	flag.IntVar(&config.bufferSize, "size", 1024, "The size of the buffer")
	flag.StringVar(&config.dispatcherType, "dispatcher", "echo", "Dispatcher type (echo, kinesis)")
	helpRequested := flag.Bool("help", false, "Print usage help and exit")

	if len(os.Args) < 2 {
		*helpRequested = true
	}

	flag.Parse()

	if *helpRequested {
		flag.Usage()
		os.Exit(0)
	}
}

func (config *Config) Validate() bool {
	switch {
	case !validateSocketType(config.socketType):
		return false
	case !validateStreamName(config.streamName, config.dispatcherType):
		return false
	default:
		return true
	}
}

func validateSocketType(socketType string) bool {
	return ValidSocketTypes[socketType]
}

func validateStreamName(streamName string, dispatcherType string) bool {
	if dispatcherType == "kinesis" {
		return len(streamName) > 0
	}
	return true
}
