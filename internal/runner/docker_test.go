package runner

import (
	"context"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/rknizzle/faas/internal/runner/mocks"
	"github.com/stretchr/testify/mock"
)

func TestRunContainer(t *testing.T) {
	mockDockerClient := new(mocks.DockerClient)
	t.Run("success", func(t *testing.T) {

		mockDockerClient.On("ContainerCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, "").Return(container.ContainerCreateCreatedBody{ID: "123"}, nil).Once()
		mockDockerClient.On("ContainerStart", mock.Anything, "123", mock.Anything).Return(nil).Once()

		dr := DockerRunner{mockDockerClient, "xxx"}
		id, err := dr.RunContainer("name")
		if err != nil {
			t.Fatalf("err %s", err)
		}

		if id != "123" {
			t.Fatalf("Returned ID did not match ID generated from ContainerCreate. Expected 123 and got %s", id)
		}
	})
}

func TestPullImage(t *testing.T) {
	mockDockerClient := new(mocks.DockerClient)
	t.Run("success", func(t *testing.T) {
		mockDockerClient.On("ImagePull", mock.Anything, "name", mock.Anything).Return(ioutil.NopCloser(strings.NewReader("test")), nil).Once()

		dr := DockerRunner{mockDockerClient, "xxx"}
		err := dr.PullImage("name")
		if err != nil {
			t.Fatalf("err %s", err)
		}
	})
}

func TestContainerIP(t *testing.T) {
	mockDockerClient := new(mocks.DockerClient)
	t.Run("success", func(t *testing.T) {
		mockDockerClient.On("ContainerInspect", mock.Anything, "xxx").Return(types.ContainerJSON{NetworkSettings: &types.NetworkSettings{DefaultNetworkSettings: types.DefaultNetworkSettings{IPAddress: "1.1.1.1"}}}, nil).Once()

		dr := DockerRunner{mockDockerClient, "xxx"}

		ctx := context.Background()
		ip, err := dr.ContainerIP(ctx, "xxx")
		if err != nil {
			t.Fatalf("err %s", err)
		}
		if ip != "1.1.1.1" {
			t.Fatalf("Returned incorrect IP. Expected 1.1.1.1, got %s", ip)
		}
	})
}
