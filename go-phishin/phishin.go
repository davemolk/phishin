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
		url := c.FormatURL("eras") 
		if c.Query != "" {
			if err := c.getAndPrintEra(ctx, url); err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("era details failure: %w", err))
				return 1
			}
			return 0
		}
		if err := c.getAndPrintEras(ctx, url); err != nil {
			fmt.Fprintln(os.Stderr, "eras list failure: %w", err)
			return 1
		}
		return 0
	case "years":
		url := c.FormatURL("years") 
		if c.Query != "" {
			if err := c.getAndPrintYear(ctx, url); err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("year details failure: %w", err))
				return 1
			}
			return 0
		}
		if err := c.getAndPrintYears(ctx, url); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("years list failure: %w", err))
			return 1
		}
		return 0
	case "songs":
		url := c.FormatURL("songs") 
		if err := c.Get(ctx, url, nil); err != nil {
			return 1
		}
	case "tours":
		url := c.FormatURL("tours") 
		if err := c.Get(ctx, url, nil); err != nil {
			return 1
		}
	case "venues":
		url := c.FormatURL("venues") 
		if err := c.Get(ctx, url, nil); err != nil {
			return 1
		}
	case "shows":
		url := c.FormatURL("shows") 
		if err := c.getAndPrintShows(ctx, url); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("shows list failure: %w", err))
			return 1
		}
		return 0
	case "show-on-date":
		url := c.FormatURL("show-on-date") 
		if err := c.Get(ctx, url, nil); err != nil {
			return 1
		}
	case "shows-on-day-of-year":
		url := c.FormatURL("shows-on-day-of-year") 
		if err := c.Get(ctx, url, nil); err != nil {
			return 1
		}
	case "random-show":
		url := c.FormatURL("random-show") 
		if err := c.Get(ctx, url, nil); err != nil {
			return 1
		}
	case "tracks":
		url := c.FormatURL("tracks") 
		if err := c.Get(ctx, url, nil); err != nil {
			return 1
		}
	case "search":
		url := c.FormatURL("search") 
		if err := c.Get(ctx, url, nil); err != nil {
			return 1
		}
	// case "playlists":
		
	case "tags":
		url := c.FormatURL("tags") 
		if err := c.Get(ctx, url, nil); err != nil {
			return 1
		}
	}

	return 0
}
