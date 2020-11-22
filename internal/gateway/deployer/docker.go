package deployer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"io/ioutil"
	"time"
)

type DockerDeployer struct {
	cli  *client.Client
	auth string
}

func NewDockerDeployer(registryUsername string, registryPassword string) (*DockerDeployer, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.40"))
	if err != nil {
		return nil, err
	}

	auth, err := generateRegistryAuth(registryUsername, registryPassword)
	if err != nil {
		return nil, err
	}

	return &DockerDeployer{cli, auth}, nil
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

// Builds a Docker image
func (d DockerDeployer) BuildImage(directory string, tag string) error {
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
func (d DockerDeployer) PushImage(image string) error {
	ctx := context.Background()
	out, err := d.cli.ImagePush(ctx, image, types.ImagePushOptions{
		RegistryAuth: d.auth,
	})
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}
