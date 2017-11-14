package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Config struct {
	socketType    string
	socketAddress string
	bufferSize    int
	echoMode      bool
}

func NewConfig() *Config {
	// parse arguments
	socketType := flag.String("type", "tcp", "The socket's type (tcp or unix)")
	socketAddress := flag.String("address", ":8888", "The socket's address (port or file)")
	bufferSize := flag.Int("size", 1024, "The size of the buffer")
	echoMode := flag.Bool("echo", false, "Activate echo mode (no kinesis requests)")
	helpRequested := flag.Bool("help", false, "Print usage help and exit")

	flag.Parse()

	if *helpRequested {
		flag.Usage()
		os.Exit(0)
	}

	// validate
	if *socketType != "tcp" && *socketType != "unix" {
		log.Fatal(fmt.Sprintf("Invalid socket type (%s)", *socketType))
	}

	return &Config{
		socketType:    *socketType,
		socketAddress: *socketAddress,
		bufferSize:    *bufferSize,
		echoMode:      *echoMode,
	}
}
