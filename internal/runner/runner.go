package runner

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
)

type Runner struct {
	CR ContainerRunner
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

// SendRequestToContainer sends an HTTP request to the containers IP and tells the container to
// start running the function code
func (r Runner) SendRequestToContainer(url string, input []byte) (string, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(input))
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
