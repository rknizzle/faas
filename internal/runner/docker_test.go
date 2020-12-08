package runner

import (
	"testing"

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
