# Instructions

This is under construction. Soon<sup>TM</sup>

## Server

build: `go build`

Usage: `./kinesis-dispatcher [type]`, where type is either `tcp` or
`unix`

* `tcp`: listens to port `localhost:8888`
* `unix`: listens to socket `/tmp/kinesis-dispatcher.sock`

## Client

usage: `ruby client.rb [type]`, where type is either `tcp` or
`unix`
