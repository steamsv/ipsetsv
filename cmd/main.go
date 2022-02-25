package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:  "ipsetsv",
	Usage: "sync ipset from db to unblock server",
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
