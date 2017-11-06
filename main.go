package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func echo(conn io.Reader) {
	io.Copy(os.Stdout, conn)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please specify unix or tcp")
		os.Exit(1)
	}

	var listener net.Listener
	var err error
	switch os.Args[1] {
	case "tcp":
		listener, err = net.Listen("tcp", ":8888")
	case "unix":
		os.Remove("/tmp/kinesis-dispatcher.sock")
		listener, err = net.Listen("unix", "/tmp/kinesis-dispatcher.sock")
	default:
		fmt.Println("Invalid type specified.")
		os.Exit(1)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(conn net.Conn, handler func(io.Reader)) {
			fmt.Println("Handling connection")
			handler(conn)
			fmt.Println("Closing connection")
			conn.Close()
		}(conn, echo)
	}
}
