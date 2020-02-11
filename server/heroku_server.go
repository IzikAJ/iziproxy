package server

// NewHerokuServer - create new Heroku Server with confguration
func NewHerokuServer(params *Config) (server *Server) {
	server = &Server{
		Host:   params.Host,
		Port:   params.Port,
		Single: params.Single,
		Stats:  Stats{},

		pool:        make(ProxyPackMap),
		spaceSignal: make(SpaceSignal),
	}

	server.tcp = NewHerokuTCPServer(server)
	server.web = NewHerokuWEBServer(server)
	// server.web = NewWEBServer(server)
	return
}
