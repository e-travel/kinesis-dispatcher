package main

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/e-travel/message-dispatcher/dispatchers"
	"github.com/e-travel/message-dispatcher/servers"
)

func createDispatcher(config *Config) (*dispatchers.MessageDispatcher, error) {
	// choose the backend
	var err error
	var svc dispatchers.Service
	switch config.dispatcherType {
	case "echo":
		svc = &dispatchers.EchoService{}
	case "influx":
		svc = dispatchers.NewInfluxService(config.influxHost, config.influxDatabase)
	case "kinesis":
		svc = dispatchers.NewKinesisService(config.streamName, config.awsRegion)
	default:
		err = errors.New(fmt.Sprintf("Invalid backend type: %s",
			config.dispatcherType))
	}
	return dispatchers.NewMessageDispatcher(svc, config.bufferSize), err
}

func main() {
	config := &Config{}
	ParseFromCommandLine(config)
	if !config.Validate() {
		log.Fatalf("Invalid socket type (%s)", config.socketType)
	}
	// create the dispatcher (backend)
	dispatcher, err := createDispatcher(config)
	if err != nil {
		log.Fatal(err.Error())
	}
	// start the worker
	go dispatcher.Dispatch()
	// setup a hook which will fire when the server is up and running
	// currently used only in tests
	running := make(chan bool)
	go func() {
		<-running
	}()
	// TODO: capture interrupt signals and stop server OR use Context
	server := servers.CreateServer(config.socketType, config.socketAddress)
	server.Serve(dispatcher, running)
}
