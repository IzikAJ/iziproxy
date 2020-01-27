package client

import (
	"net/http"
	"sync"

	"github.com/izikaj/iziproxy/shared"
)

// Client - client instance
type Client struct {
	Getaway  string
	Host     string
	Space    string
	Fallback bool

	conn   *shared.Connection
	wg     sync.WaitGroup
	http   *http.Client
	retry  int
	alive  bool
	signal chan error
}

// NewClient - create new client with confguration
func NewClient(params Config) *Client {
	return &Client{
		Getaway:  "127.0.0.1:2010",
		Host:     params.Addr,
		Space:    params.Space,
		Fallback: params.Space == "",
	}
}
