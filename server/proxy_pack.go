package server

import (
	"github.com/izikaj/iziproxy/shared"
)

// ProxyPack - just one server req
type ProxyPack struct {
	Request  shared.Request
	Response shared.Request
	signal   chan int
}
