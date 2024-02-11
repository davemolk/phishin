package cli

import (
	"fmt"
	"os"
)

// todo finish
const usage = `Usage of phishin:

phishin [command] <flags>...`


func Run() int {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		return 1
	}
	apiKey := os.Getenv("PHISHIN_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "please set the PHISHIN_API_KEY environment variable and try again")
		return 1
	}
	c := NewClient(apiKey)

	if err := c.fromArgs(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("unable to parse args: %w", err))
		return 1
	}

	endpoint := os.Args[1]
	switch endpoint {
		
	}

	return 0
}