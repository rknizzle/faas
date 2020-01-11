package main

import (
	"fmt"
	"github.com/docker/docker/client"
)

// Running core Docker functionality from Go
func main() {
	fmt.Println("Generating a client...")
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	fmt.Println(cli)
}
