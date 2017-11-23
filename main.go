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
	if config.echoMode {
		recipient = &dispatchers.Echo{}
	} else {
		recipient = dispatchers.NewKinesis(config.streamName)
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
