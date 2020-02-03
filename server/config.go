package server

import (
	"github.com/google/uuid"
)

// SpaceSignal - channel to inform about recived request
type SpaceSignal = chan uuid.UUID

// UfoSignal - channel to inform about recived error
type UfoSignal = chan error

// SpaceSignalMap - map of SpaceSignal channels
type SpaceSignalMap = map[string](SpaceSignal)

// ProxyPackMap - map of ProxyPacks
type ProxyPackMap = map[uuid.UUID]*ProxyPack

// Config - server config
type Config struct {
	Host string
	Port int

	Single bool
}
