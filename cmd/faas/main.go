package main

import (
	"fmt"
	"github.com/rknizzle/faas/api"
	"github.com/rknizzle/faas/client"
	"log"
	"os"
	"path/filepath"
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
		err := client.Invoke(os.Args[2])
		if err != nil {
			log.Fatalf(err.Error())
		}
	} else if subcommand == "start" {
		// start the server
		api.Start()
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
	fmt.Println("  init")
	fmt.Println("  build")
	fmt.Println("  invoke <function>")
}
