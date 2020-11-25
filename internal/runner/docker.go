package runner

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"io"
	"os"
)

// DockerRunner implements ContainerRunner and uses the Docker SDK to pull images and run function
// code in a Docker container
type DockerRunner struct {
	cli  *client.Client
	auth string
}

// NewDockerRunner initializes a DockerRunner with a Docker API version and credentials to access a Dockerhub registry
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

// generateRegistryAuth converts a Dockerhub username and password into an authentication string
// used by the Docker SDK
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

// PullImage downloads a container image from Dockerhub
func (d DockerRunner) PullImage(name string) error {
	ctx := context.Background()
	out, err := d.cli.ImagePull(ctx, name, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	// block until the image is fully downloaded
	// TODO: Theres probably a better way to do this
	io.Copy(os.Stdout, out)
	return nil
}

// RunContainer creates and starts a container from a local image
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

	// for now log the containers logs to stdout
	d.logOutputToConsole(ctx, resp.ID)
	return nil
}

// ContainerIP returns the IP address of a running Docker container
func (d DockerRunner) ContainerIP(ctx context.Context, id string) (string, error) {
	co, err := d.cli.ContainerInspect(ctx, id)
	if err != nil {
		return "", err
	}

	return co.NetworkSettings.IPAddress, nil
}

// output the container logs to the console
func (d DockerRunner) logOutputToConsole(ctx context.Context, id string) error {
	statusCh, errCh := d.cli.ContainerWait(ctx, id, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	out, err := d.cli.ContainerLogs(ctx, id, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return nil
}
