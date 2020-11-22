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

// Builds a Docker image
func (d DockerRunner) BuildImage(directory string, tag string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(300)*time.Second)
	defer cancel()

	// Get a tar file from directory
	dockerfileTarReader, err := archive.TarWithOptions(directory, &archive.TarOptions{})
	if err != nil {
		return err
	}

	resp, err := d.cli.ImageBuild(
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
func (d DockerRunner) PushImage(image string) error {
	ctx := context.Background()
	fmt.Println("Going to push " + image)
	out, err := d.cli.ImagePush(ctx, image, types.ImagePushOptions{
		RegistryAuth: d.auth,
	})
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
	return nil
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
	io.Copy(os.Stdout, out)
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

	statusCh, errCh := d.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := d.cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return nil
}

func (d DockerRunner) ContainerIP(ctx context.Context, id string) (string, error) {
	co, err := d.cli.ContainerInspect(ctx, id)
	if err != nil {
		return "", err
	}

	return co.NetworkSettings.IPAddress, nil
}
