package runner

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Runner struct {
	CR ContainerRunner
}

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

func (r Runner) SendRequestToContainer(ip string, input []byte) (string, error) {
	url := fmt.Sprintf("http://%s:8080/invoke", ip)

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
