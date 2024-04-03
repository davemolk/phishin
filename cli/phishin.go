package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
)

const usage = `usage: phishin <endpoint argument> [<flags>]
    	phishin is a cli client for https://phish.in/ (see https://phish.in/api-docs for more details).

request the 3 most recent shows like this:
	phishin shows -pp 3 -p 1 --sort-attr date --sort-dir desc

outputs the following:
	Date:       Venue:                 Location:                     Duration:
	2024-02-20  Moon Palace            Quintana Roo, Cancun, Mexico  54m 19s
	2023-12-31  Madison Square Garden  New York, NY                  4h 6m
	2023-12-30  Madison Square Garden  New York, NY                  2h 53m

	Total Entries: 1760  Total Pages: 587  Result Page: 1

getting started:
	get an api key (info at https://phish.in/contact-info).
	set it as an environment variable (PHISHIN_API_KEY).
	go phishin!

supported arguments:
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

arguments correspond to the phishin endpoints, and one (and only one) argument must be specified.
most allow an optional search query (-s/--search) to change the output from a list of
entities to details about a particular entity. behavior can be customized further via flags.

note: the two exceptions to the above are 'phishin help'/'phishin h' and 'phishin endpoint'/
'phishin e'.

general flags:
-s/--search		search query, format depends on the specific endpoint
--debug			print the url that is being sent to the phishin server

list-related flags:
-d/--sort-dir		direction to sort in. options are asc or desc
-a/--sort-attr		attribute to sort on (e.g. name, date)
-pp/--per-page		number of results to list per page (default is 20)
-p/--page		which page of results to display (default is 1)
-t/--tag		filter results by a specific tag (applicable for /tracks and /shows)

note: list-related flags are supported for /shows, /songs, /tracks, and /venues. they will
be ignored if you include them for other commands.

output-related flags:
-o/--output		options are json or text, default to text
-v/--verbose 		include extra information in output (not supported in all routes)

get a blank space where results should be? try the following:
format dates as "1995-12-31"
search for venues via name/past name or location ("msg" or "new york")
enter all or part of song names, tour names, etc (like "summer", "1995", "sbd", etc.)

see https://phish.in/api-docs for more details
`

const searchTips = `
get a blank space where results should be? try the following:
format dates as "1995-12-31"
search for venues via name/past name or location ("msg" or "new york")
enter all or part of song names, tour names, etc (like "summer", "1995", "sbd", etc.)

see https://phish.in/api-docs for more details
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
	erasPath           = "eras"
	yearsPath          = "years"
	songsPath          = "songs"
	toursPath          = "tours"
	venuesPath         = "venues"
	showsPath          = "shows"
	showOnDatePath     = "show-on-date"
	showsDayOfYearPath = "shows-on-day-of-year"
	randomShowPath     = "random-show"
	tracksPath         = "tracks"
	searchPath         = "search"
	tagsPath           = "tags"
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

	// get context at this point?
	// customize? or not...
	c.ErrGroup.SetLimit(4)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	path := args[0]
	if err := c.run(ctx, path); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err := c.ErrGroup.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
