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
		//client.Build()
	} else if subcommand == "invoke" {
		//client.Invoke()
	} else if subcommand == "-h" {
		helpInfo()
	} else {
		helpInfo()
	}
}

func helpInfo() {
	fmt.Println("HELP INFO HERE")
}
