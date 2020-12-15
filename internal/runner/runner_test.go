package runner

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/rknizzle/faas/internal/runner/mocks"
	"github.com/stretchr/testify/mock"
)

func TestSendRequestToContainer(t *testing.T) {
	mockHTTPPoster := new(mocks.HTTPPoster)
	t.Run("success", func(t *testing.T) {

		input := []byte("input")
		mockHTTPPoster.On(
			"Post",
			mock.Anything,
			"application/json",
			mock.Anything,
		).Return(
			&http.Response{Body: ioutil.NopCloser(strings.NewReader("test"))},
			nil,
		).Once()

		r := NewRunner(DockerRunner{}, mockHTTPPoster)
		_, err := r.SendRequestToContainer("http://1.1.1.1:8080/invoke", input)
		if err != nil {
			t.Fatalf("err %s", err)
		}
	})
}
