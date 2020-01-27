package client

import (
	"fmt"
	"net/http"
	"sync"
)

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
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	for client.alive = true; client.alive; {
		client.connect()
		client.wg.Wait()
	}
}
