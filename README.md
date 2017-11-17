# Kinesis Dispatcher

Project under construction. Soon<sup>TM</sup>.

## Purpose

The purpose of this program is to act as a server for receiving
application data and forwarding them to kinesis. Two kind of listeners
are supported: unix domain sockets and tcp sockets. Emphasis is given
to performance and non-blocking client requests.

## Instructions

### Server

Usage instructions: `./kinesis-dispatcher -help`

### Client

usage: `ruby client.rb [type]`, where type is either `tcp` or
`unix`

### Development

1. Install the [go dep tool](https://github.com/golang/dep)
2. Download the code using go get
3. `dep ensure`
4. `go build -v`
5. `go test -v`

# License

This program is available as open source under the terms of the [MIT
License](http://opensource.org/licenses/MIT).
