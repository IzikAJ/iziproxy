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

func (server *Server) place(pack *ProxyPack) {
	server.Lock()
	defer server.Unlock()
	server.pool[pack.Request.ID] = pack
}
