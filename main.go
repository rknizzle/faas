package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Running core Docker functionality from Go
func main() {
	fmt.Println("Generating a client...")
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	// list containers
	fmt.Println("Listing running containers if any...")
	listContainers(*cli)
}

// Lists all running Docker containers on the system
func listContainers(cli client.Client) error {
	// Get list of running containers
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return err
	}

	// Print out all running containers
	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
	return nil
}
