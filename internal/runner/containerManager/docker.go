package manager

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type DockerRunner struct {
	cli  *client.Client
	auth string
}

func NewDockerRunner(registryUsername string, registryPassword string) (*DockerRunner, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.40"))
	if err != nil {
		return nil, err
	}

	auth, err := generateRegistryAuth(registryUsername, registryPassword)
	if err != nil {
		return nil, err
	}

	return &DockerRunner{cli, auth}, nil
}

// Generate the dockerhub credentials from ENV variables
func generateRegistryAuth(username string, password string) (string, error) {
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
func (d DockerRunner) PullImage(name string) error {
	ctx := context.Background()
	out, err := d.cli.ImagePull(ctx, name, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Create and start a container from a local image
func (d *DockerRunner) RunContainer(image string) error {
	ctx := context.Background()

	resp, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		ExposedPorts: nat.PortSet{
			"8080/tcp": struct{}{},
		},
	}, nil, nil, "")
	if err != nil {
		return err
	}

	err = d.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (d DockerRunner) ContainerIP(ctx context.Context, id string) (string, error) {
	co, err := d.cli.ContainerInspect(ctx, id)
	if err != nil {
		return "", err
	}

	return co.NetworkSettings.IPAddress, nil
}
