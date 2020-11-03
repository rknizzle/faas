package manager

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"io"
	"io/ioutil"
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

	// block until the image is finished building
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
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

// Pull an image from a remote repository
func (m *Manager) PullImage(name string) error {
	ctx := context.Background()
	out, err := m.cli.ImagePull(ctx, name, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
	return nil
}

// Create and start a container from a local image
func (m *Manager) RunContainer(image string) error {
	ctx := context.Background()

	resp, err := m.cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		ExposedPorts: nat.PortSet{
			"8080/tcp": struct{}{},
		},
		// if Cmd is left out it will run the command specified in the Dockerfile
		//Cmd: []string{"node", "main.js"},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"8080/tcp": []nat.PortBinding{{
				HostIP:   "0.0.0.0",
				HostPort: "8080",
			}},
		},
	}, nil, "")
	if err != nil {
		return err
	}

	err = m.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	statusCh, errCh := m.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := m.cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return nil
}
