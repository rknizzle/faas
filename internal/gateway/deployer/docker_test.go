package deployer

import (
	"testing"
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
