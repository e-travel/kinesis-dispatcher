package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/e-travel/message-dispatcher/dispatchers"
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
		sess := session.Must(session.NewSession(&aws.Config{
			Retryer: client.DefaultRetryer{NumMaxRetries: 10},
			Region:  aws.String("eu-west-1"),
		}))
		recipient = &dispatchers.Kinesis{
			Session: kinesis.New(sess),
		}
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
	Serve(config, RequestHandler, buffer, running)
}
