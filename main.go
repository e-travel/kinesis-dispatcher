package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

func main() {
	config := &Config{}
	ParseFromCommandLine(config)
	if !config.Validate() {
		log.Fatal(fmt.Sprintf("Invalid socket type (%s)", config.socketType))
	}
	var dispatcher Dispatcher
	if config.echoMode {
		fmt.Println("Mode: echo")
		dispatcher = &EchoDispatcher{}
	} else {
		fmt.Println("Mode: kinesis")
		sess := session.Must(session.NewSession(&aws.Config{
			Retryer: client.DefaultRetryer{NumMaxRetries: 10},
			Region:  aws.String("eu-west-1"),
		}))
		client := &KinesisClient{
			session: kinesis.New(sess),
		}
		dispatcher = NewBufferedDispatcher(config, client)
		go dispatcher.Dispatch()
	}
	// setup a hook which will fire when the server is up and running
	// currently used only in tests
	running := make(chan bool)
	go func() {
		<-running
	}()
	// TODO: capture interrupt signals and stop server OR use Context
	Serve(config, RequestHandler, dispatcher, running)
}
