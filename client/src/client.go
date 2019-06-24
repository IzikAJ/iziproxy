package client

import (
	"net/http"
	"sync"

	"shared"
)

// Client - client instance
type Client struct {
	Getaway string
	Host    string

	conn  *shared.Connection
	wg    sync.WaitGroup
	http  *http.Client
	retry int
}
