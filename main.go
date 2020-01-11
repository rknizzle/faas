package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"io"
	"log"
	"os"
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

	fmt.Println("Generating DockerHub credentials...")
	// get dockerhub credentials
	auth, err := generateAuth()
	if err != nil {
		panic(err)
	}

	fmt.Println("Pushing the image...")
	// push
	err = pushImage(*cli, tagName, auth)
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

// Push an image to remote repository
func pushImage(cli client.Client, image string, authString string) error {
	ctx := context.Background()
	fmt.Println("Going to push " + image)
	out, err := cli.ImagePush(ctx, image, types.ImagePushOptions{
		RegistryAuth: authString,
	})
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
	return nil
}

// Generate the dockerhub credentials from ENV variables
func generateAuth() (string, error) {
	authConfig := types.AuthConfig{
		Username: os.Getenv("DOCKER_USERNAME"),
		Password: os.Getenv("DOCKER_PASSWORD"),
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	return authStr, nil
}
