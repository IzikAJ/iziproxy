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

// Start - start server daemon
func (server *Server) Start() {
	fmt.Println("Starting FULL Server...")
	defer fmt.Println("Server exists")

	server.locker.Add(2)

	// start tcp server
	go NewTCPServer(server).Start()
	// go TCPServer(config)
	// start web server
	go NewWEBServer(server).Start()
	// web := Web{}
	// go web.start(config)

	server.locker.Wait()
}

func (server *Server) placePack(pack *ProxyPack) {
	server.Lock()
	defer server.Unlock()
	server.pool[pack.Request.ID] = pack
}
