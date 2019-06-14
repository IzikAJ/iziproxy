package client

import (
	"fmt"
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

var once sync.Once

// Init will initialize client instance
func (client *Client) Init() {
	once.Do(func() {
		client.http = &http.Client{}
		(*client).retry = 10
	})
}

// Start will boot up client
func (client *Client) Start() {
	fmt.Println("Starting client...\nServe:", client.Host)
	defer fmt.Println("Client closed!")

	client.connect()

	client.wg.Wait()
}
