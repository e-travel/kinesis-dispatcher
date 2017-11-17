package main

import "fmt"

func main() {
	config := NewConfig()
	var dispatcher Dispatcher
	if config.echoMode {
		fmt.Println("Mode: echo")
		dispatcher = &EchoDispatcher{}
	} else {
		fmt.Println("Mode: kinesis")
		dispatcher = NewKinesisDispatcher(config)
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
