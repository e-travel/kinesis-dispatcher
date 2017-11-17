package main

import (
	"flag"
	"os"
)

type Config struct {
	socketType    string
	socketAddress string
	bufferSize    int
	echoMode      bool
}

func ParseFromCommandLine(config *Config) {
	flag.StringVar(&config.socketType, "type", "tcp", "The socket's type (tcp or unix)")
	flag.StringVar(&config.socketAddress, "address", ":8888", "The socket's address (port or file)")
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
	// validate
	if config.socketType != "tcp" && config.socketType != "unix" {
		return false
	}
	return true
}
