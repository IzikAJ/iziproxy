package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/izikaj/iziproxy/client"
)

func main() {
	fmt.Println("Starting client...", flag.Args())
	defer fmt.Println("THE END!")
	runtime.GOMAXPROCS(4)

	params := client.Config{
		Addr:  "http://localhost:3001",
		Space: "",
	}

	flag.StringVar(&(params.Addr), "addr", "http://localhost:3001", "proxy addr")
	flag.StringVar(&(params.Space), "space", "", "proxy space")

	flag.Parse()

	fmt.Println("params", params)
	client.NewClient(params).Init().Start()
}
