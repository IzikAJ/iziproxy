package server

import (
	"sync"

	"github.com/google/uuid"
)

// Config - server config
type Config struct {
	Host string
	Port int

	Stats  Stats
	locker sync.WaitGroup

	Single      bool
	initialized bool
	sync.Mutex
	pool  map[uuid.UUID]*ProxyPack
	space map[string](chan<- uuid.UUID)
}
