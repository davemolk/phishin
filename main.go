package main

import (
	"os"

	phishin "github.com/davemolk/phishin/go-phishin"
)

func main() {
	os.Exit(phishin.Run())
}