package server

import (
	"sync"

	"github.com/google/uuid"
)

// Cable - TODO
type Cable struct {
	Scope      string
	Owner      string `default:"test"`
	Connected  bool
	Authorized bool

	Stats       Stats
	spaceSignal chan uuid.UUID
	ufoSignal   chan error

	sync.Mutex
	pool map[uuid.UUID]*ProxyPack
}
