package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/rknizzle/faas/client"
	"github.com/rknizzle/faas/internal/gateway/api"
	"github.com/rknizzle/faas/internal/gateway/datastore"
	"github.com/rknizzle/faas/internal/gateway/deployer"
	"github.com/rknizzle/faas/internal/gateway/loadbalancer"
	"github.com/rknizzle/faas/internal/runner"
	"github.com/spf13/afero"
)

func main() {
	if len(os.Args) <= 1 {
		helpInfo()
		return
	}

	subcommand := os.Args[1]
	if subcommand == "init" {
		err := client.Init()
		if err != nil {
			log.Fatalf(err.Error())
		}
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf(err.Error())
		}

		dirName := filepath.Base(wd)
		fmt.Printf("Function %s successfully initialized.\n", dirName)
		fmt.Println("Get started by editing index.js")
	} else if subcommand == "build" {
		// verify that there is a fn.yaml file
		if _, err := os.Stat("fn.yaml"); os.IsNotExist(err) {
			fmt.Println("Not a proper directory. Run faas init to start a new function.")
			return
		}
		invokeName, err := client.Build()
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Println("invoke name:", invokeName)
	} else if subcommand == "invoke" {
		response, err := client.Invoke(os.Args[2])
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf("Response: %s", response)
	} else if subcommand == "start" {
		// start the server
		startGatewayAPI()
	} else if subcommand == "-h" {
		helpInfo()
	} else if subcommand == "--help" {
		helpInfo()
	} else {
		helpInfo()
	}
}

func helpInfo() {
	fmt.Println("usage: faas <command>")
	fmt.Println("commands:")
	fmt.Println("  start (Starts the faas server")
	fmt.Println("  init (only nodejs functions currently supported)")
	fmt.Println("  build")
	fmt.Println("  invoke <function>")
}

func startGatewayAPI() {
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

	cRunner, err := runner.NewDockerRunner(uname, password)
	if err != nil {
		fmt.Println("Failed to initialize docker for deployments")
		os.Exit(0)
	}

	run := runner.Runner{CR: cRunner}

	fs := afero.NewOsFs()
	d := deployer.NewDeployer(cDeployer, fs)
	lb := loadbalancer.NewLoadBalancer(run)
	ds := datastore.Datastore{}

	api.NewGatewayHandler(r, d, lb, ds)
	r.Run()
}
