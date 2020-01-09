package server

import (
	"fmt"

	"github.com/izikaj/iziproxy/shared"
)

// ProxyPack - just one server req
type ProxyPack struct {
	Request  shared.Request
	Response shared.Request
	signal   chan int
}

// Start - full server daemon
func Start(config *Config) {
	(*config).Initialize()

	fmt.Println("TODO: Server daemon")

	fmt.Println("Starting...")
	defer fmt.Println("THE END!")

	(*config).locker.Add(2)

	go TCPServer(config)
	web := Web{}
	go web.start(config)

	(*config).locker.Wait()
}
