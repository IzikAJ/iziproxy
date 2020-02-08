package server

import (
	"fmt"
	"sync"

	"github.com/izikaj/iziproxy/shared"
)

// AbstractServer - simplest web/tcp server interface
type AbstractServer interface {
	Start()
}

// AbstractTCPCommands - simplest tcp server commands interface
type AbstractTCPCommands interface {
	onSetup(conn *shared.Connection, cable *Cable, data shared.ConnectionSetup) error
	onResponse(conn *shared.Connection, cable *Cable, data shared.Request) error
	onPing(conn *shared.Connection, cable *Cable) error
	onPong(conn *shared.Connection, cable *Cable) error
	onUnrecognized(conn *shared.Connection, cable *Cable, data shared.Message) error
}

// Server - server instance
type Server struct {
	Host string
	Port int

	Stats  Stats
	Single bool
	locker sync.WaitGroup

	tcp AbstractServer
	web AbstractServer

	sync.Mutex
	pool  ProxyPackMap
	space SpaceSignalMap

	spaceSignal SpaceSignal
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
		return server.spaceSignal, nil
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

		spaceSignal: make(SpaceSignal),
	}

	server.tcp = NewTCPServer(server)
	server.web = NewWEBServer(server)

	return
}
