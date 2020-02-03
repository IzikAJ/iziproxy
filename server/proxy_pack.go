package server

import (
	"github.com/izikaj/iziproxy/shared"
)

// CodeSignal - channel to inform about recived status
type CodeSignal = chan int

// ProxyPack - just one server req
type ProxyPack struct {
	Request  shared.Request
	Response shared.Request
	signal   CodeSignal
}

func (server *Server) place(pack *ProxyPack) {
	server.Lock()
	defer server.Unlock()
	server.pool[pack.Request.ID] = pack
}
