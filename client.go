package main

type Client interface {
	Put([][]byte) bool
}
