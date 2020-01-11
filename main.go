package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"io"
	"log"
	"time"
)

// Running core Docker functionality from Go
func main() {
	fmt.Println("Generating a client...")
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	fmt.Println("Building a container...")
	directoryName := "nodejs-example"
	tagName := "rkneills/nodeexample"
	// build container
	err = buildImage(*cli, directoryName, tagName)
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

// Builds a Docker image
func buildImage(cli client.Client, directoryName string, tagName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(300)*time.Second)
	defer cancel()

	// Get a tar file from directory
	dockerfileTarReader, err := archive.TarWithOptions(directoryName, &archive.TarOptions{})
	if err != nil {
		return err
	}

	resp, err := cli.ImageBuild(
		ctx,
		dockerfileTarReader,
		types.ImageBuildOptions{
			Dockerfile: "Dockerfile",
			Tags:       []string{tagName},
			NoCache:    true,
			Remove:     true,
		})
	if err != nil {
		fmt.Println(err, " :unable to build docker image")
		return err
	}
	return writeToLog(resp.Body)
}

// Writes from the build response to the log
func writeToLog(reader io.ReadCloser) error {
	defer reader.Close()
	rd := bufio.NewReader(reader)
	for {
		n, _, err := rd.ReadLine()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		log.Println(string(n))
	}
	return nil
}
