package phishin

import (
	"context"
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
	c := NewClient(apiKey, os.Stdout)

	if err := c.fromArgs(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("unable to parse args: %w", err))
		return 1
	}

	// todo maybe stick url format here and pass url in?
	ctx := context.Background()

	path := os.Args[1]
	switch path {
	case "eras":
		if c.Query != "" {
			if err := c.getAndPrintEra(ctx); err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("eras details failure: %w", err))
				return 1
			}
			return 0
		}
		if err := c.getAndPrintEras(ctx); err != nil {
			fmt.Fprintln(os.Stderr, "eras list failure: %w", err)
			return 1
		}
		return 0
	case "years":
		if c.Query != "" {
			if err := c.getAndPrintYear(ctx); err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("years details failure: %w", err))
				return 1
			}
			return 0
		}
		if err := c.getAndPrintYears(ctx); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("years list failure: %w", err))
			return 1
		}
		return 0
	case "songs":
		if err := c.Get(ctx, "songs", nil); err != nil {
			return 1
		}
	case "tours":
		if err := c.Get(ctx, "tours", nil); err != nil {
			return 1
		}
	case "venues":
		if err := c.Get(ctx, "venues", nil); err != nil {
			return 1
		}
	case "shows":
		if err := c.getAndPrintShows(ctx); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("shows list failure: %w", err))
			return 1
		}
		return 0
	case "show-on-date":
		if err := c.Get(ctx, "show-on-date", nil); err != nil {
			return 1
		}
	case "shows-on-day-of-year":
		if err := c.Get(ctx, "shows-on-day-of-year", nil); err != nil {
			return 1
		}
	case "random-show":
		if err := c.Get(ctx, "random-show", nil); err != nil {
			return 1
		}
	case "tracks":
		if err := c.Get(ctx, "tracks", nil); err != nil {
			return 1
		}
	case "search":
		if err := c.Get(ctx, "search", nil); err != nil {
			return 1
		}
	case "playlists":
		if err := c.Get(ctx, "playlists", nil); err != nil {
			return 1
		}
	case "tags":
		if err := c.Get(ctx, "tags", nil); err != nil {
			return 1
		}
	}

	return 0
}
