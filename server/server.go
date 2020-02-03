package server

import (
	"fmt"
	"sync"
)

// Server - server instance
type Server struct {
	Host string
	Port int

	Stats  Stats
	Single bool
	locker sync.WaitGroup

	tcp *TCPServer
	web *WEBServer

	sync.Mutex
	pool  ProxyPackMap
	space SpaceSignalMap

	globalSpaceSignal SpaceSignal
}

// Start - start server daemon
func (server *Server) Start() {
	fmt.Println("Starting Server...")
	defer fmt.Println("Server exists")
	server.locker.Add(2)

	// start tcp server
	go server.tcp.Start()
	// start web server
	go server.web.Start()

	server.locker.Wait()
}

func (server *Server) findSpaceSignal(params spaceParams) (SpaceSignal, error) {
	if server.Single {
		return server.globalSpaceSignal, nil
	}
	if signal, ok := server.space[params.subdomain]; ok {
		return signal, nil
	}
	return nil, fmt.Errorf("not found")
}

// NewServer - create new Server with confguration
func NewServer(params *Config) (server *Server) {
	server = &Server{
		Host:   params.Host,
		Port:   params.Port,
		Single: params.Single,
		Stats:  Stats{},

		pool:  make(ProxyPackMap),
		space: make(SpaceSignalMap),

		globalSpaceSignal: make(SpaceSignal),
	}

	server.tcp = NewTCPServer(server)
	server.web = NewWEBServer(server)

	return
}
