package loadbalancer

import (
	"time"

	"github.com/rknizzle/faas/internal/runner"
)

type LoadBalancer struct {
	r runner.Runner
}

func NewLoadBalancer(r runner.Runner) LoadBalancer {
	return LoadBalancer{r}
}

func (lb LoadBalancer) SendToRunner(image string, input string) (string, error) {
	// start the container and return its IP address
	ip, err := lb.r.StartFnContainer(image)
	if err != nil {
		return "", err
	}

	// TODO: probably going to need a sleep here until I figure out how to know when the container web
	// server up and ready
	time.Sleep(2 * time.Second)

	// pass the user input and trigger the fn code by sending an HTTP request to the containers IP
	// address
	output, err := lb.r.SendRequestToContainer(ip, input)
	if err != nil {
		return "", err
	}

	// return the results from the fn code
	return output, nil
}
