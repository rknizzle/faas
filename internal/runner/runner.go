package runner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// ContainerRunner handles getting the function code, starting it in a container and getting the IP
// address of the container to interact with it
type ContainerRunner interface {
	PullImage(string) error
	RunContainer(string) (string, error)
	ContainerIP(context.Context, string) (string, error)
}

// Runner handles starting up the function code in a container, sending an HTTP request to invoke
// the function, and returning the response of the function back to the caller
type Runner struct {
	CR     ContainerRunner
	Client HTTPPoster
}

type HTTPPoster interface {
	Post(string, string, io.Reader) (*http.Response, error)
}

func NewRunner(cr ContainerRunner) Runner {
	return Runner{CR: cr, Client: &http.Client{}}
}

// StartFnContainer starts a container containing the function code and returns the IP address of
// the container so that the function can be invoked via HTTP request
func (r Runner) StartFnContainer(image string) (string, error) {
	id, err := r.CR.RunContainer(image)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	ip, err := r.CR.ContainerIP(ctx, id)
	if err != nil {
		return "", err
	}

	return ip, nil
}

// TriggerContainer will tell the container to run function code once the container is running and
// ready and gets the output of the function to return back to the caller
func (r Runner) TriggerContainerFn(ip string, input []byte) (output string, err error) {
	url := fmt.Sprintf("http://%s:8080/invoke", ip)

	var success bool = false
	// loop until the request to the container gets a successful response
	// TODO: add a timeout or limit to the number of requests before returning an error
	for success == false {
		output, err = r.SendRequestToContainer(url, input)
		if err == nil {
			success = true
		}

		if success == false {
			// do a small sleep after a container request fails
			time.Sleep(50 * time.Millisecond)
		}
	}

	return
}

// SendRequestToContainer sends an HTTP request to the containers IP and tells the container to
// start running the function code
func (r Runner) SendRequestToContainer(url string, input []byte) (string, error) {
	resp, err := r.Client.Post(url, "application/json", bytes.NewBuffer(input))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	bodyString := string(bodyBytes)
	return bodyString, nil
}
