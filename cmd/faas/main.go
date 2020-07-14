package main

import (
	"fmt"
	"github.com/rknizzle/faas/client"
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatalf("No arguments included")
	}

	subcommand := os.Args[1]
	if subcommand == "init" {
		err := client.Init()
		if err != nil {
			log.Fatalf(err.Error())
		}
	} else if subcommand == "build" {
		// verify that there is a fn.yaml file
		if _, err := os.Stat("fn.yaml"); os.IsNotExist(err) {
			log.Fatalf("Not a proper directory")
		}
		invokeName, err := client.Build()
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Println("invokeName:", invokeName)
	} else if subcommand == "invoke" {
		err := client.Invoke(os.Args[2])
		if err != nil {
			log.Fatalf(err.Error())
		}
	} else if subcommand == "-h" {
		helpInfo()
	} else {
		helpInfo()
	}
}

func helpInfo() {
	fmt.Println("HELP INFO HERE")
}
