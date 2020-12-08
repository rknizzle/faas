package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	gw "github.com/rknizzle/faas/internal/gateway/api"
	"github.com/rknizzle/faas/internal/gateway/datastore"
	"github.com/rknizzle/faas/internal/gateway/deployer"
	"github.com/rknizzle/faas/internal/gateway/loadbalancer"
	"os"
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

	e := deployer.NewExtractor()
	d := deployer.NewDeployer(cDeployer, e)
	lb := loadbalancer.LoadBalancer{}
	ds := datastore.Datastore{}

	gw.NewGatewayHandler(r, d, lb, ds)
	r.Run()
}
