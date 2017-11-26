package main

import (
	"fmt"
	"log"

	"github.com/e-travel/message-dispatcher/dispatchers"
	"github.com/e-travel/message-dispatcher/servers"
)

func main() {
	config := &Config{}
	ParseFromCommandLine(config)
	if !config.Validate() {
		log.Fatal(fmt.Sprintf("Invalid socket type (%s)", config.socketType))
	}
	// choose the backend
	var recipient dispatchers.Dispatcher
	switch config.dispatcherType {
	case "echo":
		recipient = &dispatchers.Echo{}
	case "kinesis":
		recipient = dispatchers.NewKinesis(config.streamName)
		go recipient.Dispatch()
	default:
		log.Fatal(fmt.Sprintf("Invalid dispatcher type (%s)", config.dispatcherType))
	}
	// create the intermediate buffer
	buffer := dispatchers.NewMessageBuffer(config.bufferSize, recipient)
	go buffer.Dispatch()
	// setup a hook which will fire when the server is up and running
	// currently used only in tests
	running := make(chan bool)
	go func() {
		<-running
	}()
	// TODO: capture interrupt signals and stop server OR use Context
	server := servers.CreateServer(config.socketType, config.socketAddress)
	server.Serve(buffer, running)
}
