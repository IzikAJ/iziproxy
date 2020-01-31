package server

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// Server - server instance
type Server struct {
	Host string
	Port int

	Stats  Stats
	Single bool
	locker sync.WaitGroup

	sync.Mutex
	pool  map[uuid.UUID]*ProxyPack
	space map[string](chan<- uuid.UUID)
}

// Start - start server daemon
func (server *Server) Start() {
	fmt.Println("Starting Server...")
	defer fmt.Println("Server exists")
	server.locker.Add(2)

	// start tcp server
	go NewTCPServer(server).Start()

	// start web server
	go NewWEBServer(server).Start()

	server.locker.Wait()
}

// NewServer - create new Server with confguration
func NewServer(params *Config) *Server {
	return &Server{
		Host:  params.Host,
		Port:  params.Port,
		Stats: Stats{},

		pool:  make(map[uuid.UUID]*ProxyPack),
		space: make(map[string](chan<- uuid.UUID)),
	}
}
