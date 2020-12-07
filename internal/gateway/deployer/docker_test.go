package deployer

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/rknizzle/faas/internal/gateway/deployer/mocks"
	"github.com/stretchr/testify/mock"
)

func TestGenerateRegistryAuth(t *testing.T) {
	expected := "eyJ1c2VybmFtZSI6ImV4YW1wbGVOYW1lIiwicGFzc3dvcmQiOiJleGFtcGxlUGFzcyJ9"

	got, err := generateRegistryAuth("exampleName", "examplePass")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if got != expected {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}

func TestGenerateRegistryAuthMissingUsername(t *testing.T) {
	_, err := generateRegistryAuth("", "examplePass")
	if err == nil {
		t.Fatal("Expected test case to fail due to missing username")
	}
}

func TestGenerateRegistryAuthMissingPassword(t *testing.T) {
	_, err := generateRegistryAuth("exampleUser", "")
	if err == nil {
		t.Fatal("Expected test case to fail due to missing password")
	}
}

func TestTestBuildImage(t *testing.T) {
	mockDockerClient := new(mocks.DockerClient)
	t.Run("success", func(t *testing.T) {

		mockDockerClient.On(
			"ImageBuild",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(
			types.ImageBuildResponse{
				Body: ioutil.NopCloser(strings.NewReader("test")),
			},
			nil,
		).Once()

		dd := DockerDeployer{mockDockerClient, "xxx"}
		err := dd.BuildImage("name", "tag")
		if err != nil {
			t.Fatalf("err %s", err)
		}
	})
}
