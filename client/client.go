package client

import (
	"net/http"
	"sync"

	"github.com/izikaj/iziproxy/shared"
)

// Client - client instance
type Client struct {
	Getaway string
	Host    string
	Space   string

	conn  *shared.Connection
	wg    sync.WaitGroup
	http  *http.Client
	retry int
}
