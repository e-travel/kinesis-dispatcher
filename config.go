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
	influxHost     string
	influxDatabase string
	bufferSize     int
	dispatcherType string
}

var ValidSocketTypes = map[string]bool{
	servers.UNIXGRAM: true,
	servers.UDP:      true,
}

func ParseFromCommandLine(config *Config) {
	flag.StringVar(&config.dispatcherType, "dispatcher", "echo", "Dispatcher type (echo, kinesis, influx)")
	flag.StringVar(&config.socketType, "type", servers.UNIXGRAM, "The socket's type (unixgram, udp)")
	flag.StringVar(&config.socketAddress, "address", "/tmp/msg-dsp.sock", "The socket's address (file)")
	flag.StringVar(&config.streamName, "stream-name", "", "The name of the kinesis stream")
	flag.StringVar(&config.awsRegion, "aws-region", "eu-west-1", "The kinesis stream's AWS region")
	flag.StringVar(&config.influxHost, "influx-host", "http://localhost:8086", "Influx server hostname")
	flag.StringVar(&config.influxDatabase, "influx-database", "", "Influx database to use")
	flag.IntVar(&config.bufferSize, "size", 1024, "The size of the buffer")
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
	case !validateInflux(config.dispatcherType, config.influxHost, config.influxDatabase):
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

func validateInflux(dispatcherType string, influxHost string, influxDatabase string) bool {
	if dispatcherType == "influx" {
		return influxHost != "" && influxDatabase != ""
	}
	return true
}
