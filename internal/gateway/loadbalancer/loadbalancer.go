package loadbalancer

import (
	"fmt"
	"github.com/rknizzle/faas/internal/runner"
	"os"
)

type LoadBalancer struct{}

func (lb LoadBalancer) SendToRunner(image string) error {
	// something like this on a runner machine. Will probably send an HTTP request to the machine
	// telling it to start running the function
	// runner.Invoke(fn)

	// This is just a hardcoded invocation of a function in a docker container for testing out the
	// flow
	// TODO: this should be replaced by an HTTP request to a runner that the LB decides to run the
	// function on
	cRunner, err := runner.NewDockerRunner(os.Getenv("DOCKER_USERNAME"), os.Getenv("DOCKER_PASSWORD"))
	if err != nil {
		return err
	}
	err = cRunner.PullImage("rkneills/fn:latest")
	if err != nil {
		return err
	}
	err = cRunner.RunContainer("rkneills/fn:latest")
	if err != nil {
		return err
	}
	return nil
}
