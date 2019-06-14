package server

import (
	"fmt"

	"shared"
)

// ProxyPack - just one server req
type ProxyPack struct {
	Request  shared.Request
	Response shared.Request
	signal   chan int
}

// Server - full server daemon
func Server(config *Config) {
	(*config).Initialize()

	fmt.Println("TODO: Server daemon")

	fmt.Println("Starting...")
	defer fmt.Println("THE END!")

	(*config).locker.Add(2)

	go TCPServer(config)
	web := Web{
		port: 1234,
		host: "0.0.0.0",
	}
	go web.start(config)

	(*config).locker.Wait()
}
