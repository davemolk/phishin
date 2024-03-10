package main

import (
	"os"

	"github.com/davemolk/phishin/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
