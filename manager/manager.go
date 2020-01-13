package manager

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"io"
	"os"
	"time"
)

type Manager struct {
	cli  *client.Client
	auth string
}

func New() *Manager {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	auth, err := generateAuth()
	if err != nil {
		panic(err)
	}

	return &Manager{cli, auth}
}

// Builds a Docker image
func (m Manager) BuildImage(directory string, tag string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(300)*time.Second)
	defer cancel()

	// Get a tar file from directory
	dockerfileTarReader, err := archive.TarWithOptions(directory, &archive.TarOptions{})
	if err != nil {
		return err
	}

	resp, err := m.cli.ImageBuild(
		ctx,
		dockerfileTarReader,
		types.ImageBuildOptions{
			Dockerfile: "Dockerfile",
			Tags:       []string{tag},
			NoCache:    true,
			Remove:     true,
		})
	if err != nil {
		fmt.Println(err, " :unable to build docker image")
		return err
	}
	io.Copy(os.Stdout, resp.Body)
	return nil
}

// Push an image to remote repository
func (m Manager) PushImage(image string) error {
	ctx := context.Background()
	fmt.Println("Going to push " + image)
	out, err := m.cli.ImagePush(ctx, image, types.ImagePushOptions{
		RegistryAuth: m.auth,
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
	username := os.Getenv("DOCKER_USERNAME")
	password := os.Getenv("DOCKER_PASSWORD")

	if username == "" || password == "" {
		return "", errors.New("Missing Dockerhub username or password")
	}

	authConfig := types.AuthConfig{
		Username: username,
		Password: password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	return authStr, nil
}
