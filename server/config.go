package server

import (
	"sync"

	"github.com/google/uuid"
)

// Config - TODO
type Config struct {
	Host string
	Port int

	Stats  Stats
	locker sync.WaitGroup

	sync.Mutex
	pool  map[uuid.UUID]*ProxyPack
	space map[string](chan<- uuid.UUID)
}

// Initialize - init config
func (conf *Config) Initialize() *Config {
	(*conf).pool = make(map[uuid.UUID]*ProxyPack)
	(*conf).space = make(map[string](chan<- uuid.UUID))
	(*conf).Stats = Stats{}
	return conf
}

// Conf - conf instance
var Conf = Config{}
