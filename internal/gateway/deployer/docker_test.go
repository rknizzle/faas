package deployer

import (
	"context"
	"io"
	"testing"

	"github.com/docker/docker/api/types"
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

// mocks for DockerClient
type fakeIOReadCloser struct{}

func (f fakeIOReadCloser) Read([]byte) (int, error) { return 0, nil }
func (f fakeIOReadCloser) Close() error             { return nil }

type mockDockerClient struct{}

func (m mockDockerClient) ImageBuild(context.Context, io.Reader, types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	return types.ImageBuildResponse{}, nil
}

func (m mockDockerClient) ImagePush(context.Context, string, types.ImagePushOptions) (io.ReadCloser, error) {
	return fakeIOReadCloser{}, nil
}

func newMockDockerDeployer() DockerDeployer {
	dc := mockDockerClient{}
	return DockerDeployer{dc, "xxx"}
}

func TestBuildImage(t *testing.T) {
	readAll = func(io.Reader) ([]byte, error) { return nil, nil }
	dd := newMockDockerDeployer()
	err := dd.BuildImage("name", "tag")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}
