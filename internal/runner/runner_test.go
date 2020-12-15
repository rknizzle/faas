package runner

import (
	"errors"
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

	t.Run("Returns an error when the HTTP request returns an error", func(t *testing.T) {

		input := []byte("input")
		mockHTTPPoster.On(
			"Post",
			mock.Anything,
			"application/json",
			mock.Anything,
		).Return(
			nil,
			errors.New("test error"),
		).Once()

		r := NewRunner(DockerRunner{}, mockHTTPPoster)
		_, err := r.SendRequestToContainer("http://1.1.1.1:8080/invoke", input)
		if err == nil {
			t.Fatal("Expected SendRequestToContainer to throw an error")
		}
	})
}

func TestTriggerContainerFn(t *testing.T) {
	mockHTTPPoster := new(mocks.HTTPPoster)
	t.Run("Function returns successfully when HTTP req returns successfully on the first try", func(t *testing.T) {

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
		_, err := r.TriggerContainerFn("1.1.1.1", input)
		if err != nil {
			t.Fatalf("err %s", err)
		}
	})

	t.Run("Function returns successfully when HTTP req requires multiple attempts", func(t *testing.T) {

		input := []byte("input")
		// first 2 attempts return an error
		mockHTTPPoster.On(
			"Post",
			mock.Anything,
			"application/json",
			mock.Anything,
		).Return(
			nil,
			errors.New("test error"),
		).Twice()

		// 3rd attempt return successfully
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
		_, err := r.TriggerContainerFn("1.1.1.1", input)
		if err != nil {
			t.Fatalf("err %s", err)
		}
	})
}
