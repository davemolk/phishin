package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// todo finish
const usage = `usage: phishin <command> [<flags>]
    phishin is a cli client for https://phish.in/ (see https://phish.in/api-docs for more details).

request the 3 most recent shows like so:
	phishin shows -pp 3 -p 1 --sort-attr date --sort-dir desc

outputs the following:
	Date:       Venue:                 Location:                     Duration:
	2024-02-20  Moon Palace            Quintana Roo, Cancun, Mexico  54m 19s
	2023-12-31  Madison Square Garden  New York, NY                  4h 6m
	2023-12-30  Madison Square Garden  New York, NY                  2h 53m

	Total Entries: 1760  Total Pages: 587  Result Page: 1

getting started:
	1) get an api key (info at https://phish.in/contact-info).
	2) set it as an environment variable (PHISHIN_API_KEY).
	3) go phishin!

commands:
	eras 			(-s as era, e.g. 3.0)
	years 			(-s as year, e.g. 1994)
	songs 			(-s as song slug or song-id, e.g. harry-hood)
	tours 			(-s as tour slug or tour id, e.g. 1983-tour)
	venues 			(-s as venue slug or venue id, e.g. the-academy)
	shows 			(-s as show date or show id, e.g. 1994-10-31)
	show-on-date -s 	(query required, format as yyyy-mm-dd)
	shows-on-day-of-year -s (query required, format as 10-31)
	random-show
	tracks 			(-s as tracks id, e.g. 6693)
	search -s
	tags 			(-s as tag slug or tag id, e.g. sbd)

most commands allow an optional search query (-s/--search) to change the output from a list 
of entities to details about a particular entity. see 'phishin endpoints' or 'phishin e' for 
details about endpoints.

general flags:
-s/--search		search query, format depends on endpoint
--debug			print the url that is being sent to the phishin server

list-related flags:
-d/--sort-dir	direction to sort in. options are asc or desc
-a/--sort-attr	attribute to sort on (e.g. name, date)
-pp/--per-page	number of results to list per page (default is 20)
-p/--page	which page of results to display (default is 1)
-t/--tag	filter results by a specific tag (applicable for /tracks and /shows)

output-related flags:
-o/--output
-v/--verbose
`

const endpointList = `

supported endpoints:

/eras
/eras/:era

/years
/years/:year

/songs
/songs/:id
/songs/:slug

/tours
/tours/:id
/tours/:slug

/venues
/venues/:id
/venues/:slug

/shows
/shows/:id
/shows/:date(yyyy-mm-dd)

/show-on-date/:date(yyyy-mm-dd)

/shows-on-day-of-year/:day(mm-dd)

/random-show

/tracks
/tracks/:id

/search/:term

/tags
/tags/:id
/tags/:slug

example usage to get era 2.0: 
phishin eras -s 2.0

see https://phish.in/api-docs for more details`

const (
	erasPath = "eras"
	yearsPath = "years"
	songsPath = "songs"
	toursPath = "tours"
	venuesPath = "venues"
	showsPath = "shows"
	showOnDatePath = "show-on-date"
	showsDayOfYearPath = "shows-on-day-of-year"
	randomShowPath = "random-show"
	tracksPath = "tracks"
	searchPath = "search"
	tagsPath = "tags"
)

func Run(args []string) int {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, usage)
		return 1
	}
	switch strings.ToLower(args[0]) {
	case "help", "h", "-help", "-h", "--help":
		fmt.Fprint(os.Stderr, usage)
		return 0
	case "endpoints", "e", "-endpoints", "-e", "--endpoints":
		fmt.Fprintln(os.Stderr, endpointList)
		return 0
	}
	apiKey := os.Getenv("PHISHIN_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "please set the PHISHIN_API_KEY environment variable and try again")
		fmt.Fprintln(os.Stderr, "keys may be requested via https://phish.in/contact-info")
		return 1
	}
	c := NewClient(apiKey, os.Stdout)

	if err := c.fromArgs(args); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("unable to parse args: %w", err))
		return 1
	}

	// todo maybe stick url format here and pass url in?
	ctx := context.Background()

	path := args[0]
	switch path {
	case erasPath:
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
	case yearsPath:
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
	case songsPath:
		url := c.FormatURL(path) 
		if c.Query != "" {
			if err := c.getAndPrintSong(ctx, url); err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("song details failure: %w", err))
				return 1
			}
			return 0
		}
		if err := c.getAndPrintSongs(ctx, url); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("songs list failure: %w", err))
			return 1
		}
		return 0
	case toursPath:
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
	case showsPath:
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
	case showOnDatePath:
		url := c.FormatURL(path) 
		if err := c.getAndPrintShow(ctx, url); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("show details failure: %w", err))
			return 1
		}
		return 0
	case showsDayOfYearPath:
		url := c.FormatURL(path)
		if err := c.getAndPrintShows(ctx, url); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("shows list failure: %w", err))
			return 1
		}
		return 0
	case randomShowPath:
		url := c.FormatURL(path)
		if err := c.getAndPrintShow(ctx, url); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("show details failure: %w", err))
			return 1
		}
		return 0
	case tracksPath:
		url := c.FormatURL(path) 
		if c.Query != "" {
			if err := c.getAndPrintTrack(ctx, url); err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("track details failure: %w", err))
				return 1
			}
			return 0
		}
		if err := c.getAndPrintTracks(ctx, url); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("tracks list failure: %w", err))
			return 1
		}
		return 0
	case searchPath:
		url := c.FormatURL(path) 
		if err := c.Get(ctx, url, nil); err != nil {
			return 1
		}
	// case "playlists":
		
	case tagsPath:
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
