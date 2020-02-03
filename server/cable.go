package server

import (
	"sync"
)

// Cable - TODO
type Cable struct {
	Scope      string
	Owner      string `default:"test"`
	Connected  bool
	Authorized bool

	Stats       Stats
	spaceSignal SpaceSignal
	ufoSignal   UfoSignal

	sync.Mutex
	pool ProxyPackMap
}
