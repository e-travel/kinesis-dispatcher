package main

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/e-travel/message-dispatcher/dispatchers"
	"github.com/e-travel/message-dispatcher/servers"
)

func createDispatcher(name string) (dispatchers.Dispatcher, error) {
	// choose the backend
	var dispatcher dispatchers.Dispatcher
	var err error
	switch name {
	case "echo":
		dispatcher = &dispatchers.Echo{}
	case "kinesis":
		dispatcher = dispatchers.NewKinesis(name)
	default:
		err = errors.New(fmt.Sprintf("Invalid dispatcher type: %s", name))
	}
	return dispatcher, err
}

func main() {
	config := &Config{}
	ParseFromCommandLine(config)
	if !config.Validate() {
		log.Fatalf("Invalid socket type (%s)", config.socketType)
	}
	// create the dispatcher (backend)
	dispatcher, err := createDispatcher(config.dispatcherType)
	if err != nil {
		log.Fatal(err.Error())
	}
	// start the worker
	go dispatcher.Dispatch()
	// create the intermediate buffer
	buffer := dispatchers.NewMessageBuffer(config.bufferSize, dispatcher)
	// start the worker
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
