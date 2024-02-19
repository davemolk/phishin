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
		url := c.FormatURL(path) 
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
		url := c.FormatURL(path) 
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
		url := c.FormatURL(path) 
		if err := c.Get(ctx, url, nil); err != nil {
			return 1
		}
	case "tours":
		url := c.FormatURL(path) 
		if c.Query != "" {
			if err := c.getAndPrintTour(ctx, url); err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("tour details failure: %w", err))
				return 1
			}
			return 0
		}
		if err := c.getAndPrintTours(ctx, url); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("tours list failure: %w", err))
			return 1
		}
		return 0
	case "venues":
		url := c.FormatURL(path) 
		if c.Query != "" {
			if err := c.getAndPrintVenue(ctx, url); err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("venue details failure: %w", err))
				return 1
			}
			return 0
		}
		if err := c.getAndPrintVenues(ctx, url); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("venues list failure: %w", err))
			return 1
		}
		return 0
	case "shows":
		url := c.FormatURL(path) 
		if c.Query != "" {
			if err := c.getAndPrintShow(ctx, url); err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("show details failure: %w", err))
				return 1
			}
			return 0
		}
		if err := c.getAndPrintShows(ctx, url); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("shows list failure: %w", err))
			return 1
		}
		return 0
	case "show-on-date":
		url := c.FormatURL(path) 
		if err := c.getAndPrintShow(ctx, url); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("show details failure: %w", err))
			return 1
		}
		return 0
	case "shows-on-day-of-year":
		url := c.FormatURL(path) 
		if err := c.Get(ctx, url, nil); err != nil {
			return 1
		}
	case "random-show":
		url := c.FormatURL(path)
		if err := c.getAndPrintShow(ctx, url); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("show details failure: %w", err))
			return 1
		}
		return 0
	case "tracks":
		url := c.FormatURL(path) 
		if err := c.Get(ctx, url, nil); err != nil {
			return 1
		}
	case "search":
		url := c.FormatURL(path) 
		if err := c.Get(ctx, url, nil); err != nil {
			return 1
		}
	// case "playlists":
		
	case "tags":
		url := c.FormatURL(path)
		if c.Query != "" {
			if err := c.getAndPrintTag(ctx, url); err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("tag details failure: %w", err))
				return 1
			}
			return 0
		}
		if err := c.getAndPrintTags(ctx, url); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("tags list failure: %w", err))
			return 1
		}
		return 0
	}

	return 0
}
