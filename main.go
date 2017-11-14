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
	Serve(config, RequestHandler, dispatcher)
}
