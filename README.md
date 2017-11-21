# Message Dispatcher

Project under construction. Soon<sup>TM</sup>.

## Purpose

The purpose of this program is to act as a server for receiving
application data and forwarding them to a recipient backend (e.g. AWS
Kinesis). Supported servers include:

* unix domain socket (stream)
* tcp socket (stream)
* unix datagram socket (datagram)

Emphasis is given to performance and non-blocking client requests.

## Instructions

### Server

Usage instructions: `./message-dispatcher -help`

### Clients

Ruby client examples can be found in directory `clients/`

### Development

0. Install `go`
1. Install the [go dep tool](https://github.com/golang/dep)
2. Download the code using `go get github.com/e-travel/message-dispatcher`
3. `dep ensure`
4. `go build -v -i`
5. `go test ./... -v`

# License

This program is available as open source under the terms of the [MIT
License](http://opensource.org/licenses/MIT).
