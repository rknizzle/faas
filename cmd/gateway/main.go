package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rknizzle/faas/internal/gateway/api"
	"github.com/rknizzle/faas/internal/gateway/datastore"
	"github.com/rknizzle/faas/internal/gateway/deployer"
	"github.com/rknizzle/faas/internal/gateway/loadbalancer"
	"github.com/spf13/afero"
)

func main() {
	uname := os.Getenv("DOCKER_USERNAME")
	password := os.Getenv("DOCKER_PASSWORD")
	if len(uname) == 0 || len(password) == 0 {
		fmt.Println("Missing Docker username or password")
		os.Exit(0)
	}

	r := gin.Default()
	cDeployer, err := deployer.NewDockerDeployer(uname, password)
	if err != nil {
		fmt.Println("Failed to initialize docker for deployments")
		os.Exit(0)
	}

	fs := afero.NewOsFs()
	d := deployer.NewDeployer(cDeployer, fs)
	lb := loadbalancer.LoadBalancer{}
	ds := datastore.Datastore{}

	api.NewGatewayHandler(r, d, lb, ds)
	r.Run()
}
