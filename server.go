package main

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func RequestHandler(conn io.Reader, dispatcher Dispatcher) {
	b, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Print("Error reading from connection")
		return
	}
	dispatcher.Put(b)
}

func Serve(config *Config, handler func(io.Reader, Dispatcher), dispatcher Dispatcher) {
	log.Printf("%s socket server listening to %s", config.socketType, config.socketAddress)
	// remove any existing socket file
	if config.socketType == "unix" {
		os.Remove(config.socketAddress)
	}
	listener, err := net.Listen(config.socketType, config.socketAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(conn net.Conn, handler func(io.Reader, Dispatcher)) {
			handler(conn, dispatcher)
			conn.Close()
		}(conn, handler)
	}
}
