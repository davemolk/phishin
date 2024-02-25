package phishin

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
)

type Client struct {
	Tabwriter *tabwriter.Writer
	HTTPClient *http.Client
	BaseURL string
	APIKey string
	PrintJSON bool
	Query string
	Parameters []string
	Output io.Writer
	Verbose bool
}

func NewClient(apiKey string, output io.Writer) *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
		BaseURL: "https://phish.in/api/v1",
		APIKey: apiKey,
		Output: output,
		Tabwriter: tabwriter.NewWriter(output, 0, 4, 2, ' ', tabwriter.DiscardEmptyColumns),
	}
}

func (c *Client) FormatURL(path string) string {
	if c.Query != "" {
		// return now to avoid mixing in params
		return fmt.Sprintf("%s/%s/%s", c.BaseURL, path, c.Query)
	}
	url := fmt.Sprintf("%s/%s", c.BaseURL, path)
	if len(c.Parameters) != 0 {
		params := strings.Join(c.Parameters, "&")
		url = fmt.Sprintf("%s?%s", url, params)
	}
	return url
}

func (c *Client) fromArgs(args []string) error {
	// todo use one flagset for all (or keep separate in case things change?)
	// separate means we can have diff flag usage msgs
	eras := flag.NewFlagSet("eras", flag.ExitOnError)
	eras.StringVar(&c.Query, "s", "", "search query")
	erasOutput := eras.String("o", "text", "print output as text/json")
	// eras.Usage = func() {
	// 	fmt.Fprintf(flag.CommandLine.Output(), "%s tool\n", os.Args[0])
	// 	fmt.Fprintln(flag.CommandLine.Output(), "usage information:")
	// 	flag.PrintDefaults()
	// }

	years := flag.NewFlagSet("years", flag.ExitOnError)
	years.StringVar(&c.Query, "s", "", "search query")
    yearsOutput := years.String("o", "text", "print output as text/json")

	tags := flag.NewFlagSet("tags", flag.ExitOnError)
	tags.StringVar(&c.Query, "s", "", "search query")
	tagsOutput := tags.String("o", "text", "print output as text/json")

	tours := flag.NewFlagSet("tours", flag.ExitOnError)
	tours.StringVar(&c.Query, "s", "", "search query")
	toursOutput := tours.String("o", "text", "print output as text/json")

	/* longer ones */

	shows := flag.NewFlagSet("shows", flag.ExitOnError)
	shows.StringVar(&c.Query, "s", "", "search query")
	showsSortDir := shows.String("sortDir", "", "sort results asc/desc") 
    showsSortAttr := shows.String("sortAttr", "", "sort results <attr>")
	showsPerPage := shows.Int("pp", 20, "per page")
	showsPage := shows.Int("p", 1, "page")
	showsTag := shows.String("tag", "", "filter by <tag>")
	showsVerbose := shows.Bool("v", false, "fill this out")
	showsOutput := shows.String("o", "text", "print output as text/json")

	venues := flag.NewFlagSet("venues", flag.ExitOnError)
	venues.StringVar(&c.Query, "s", "", "search query")
	venuesSortDir := venues.String("sortDir", "", "sort results asc/desc") 
    venuesSortAttr := venues.String("sortAttr", "", "sort results <attr>")
	venuesPerPage := venues.Int("pp", 20, "per page")
	venuesPage := venues.Int("p", 1, "page")
	venuesOutput := venues.String("o", "text", "print output as text/json")
    
	songs := flag.NewFlagSet("songs", flag.ExitOnError)
	songs.StringVar(&c.Query, "s", "", "search query")
	songsSortDir := songs.String("sortDir", "", "sort results asc/desc") 
    songsSortAttr := songs.String("sortAttr", "", "sort results <attr>")
	songsPerPage := songs.Int("pp", 20, "per page")
	songsPage := songs.Int("p", 1, "page")
	songsOutput := songs.String("o", "text", "print output as text/json")

	tracks := flag.NewFlagSet("tracks", flag.ExitOnError)
	tracks.StringVar(&c.Query, "s", "", "search query")
	tracksSortDir := tracks.String("sortDir", "", "sort results asc/desc") 
    tracksSortAttr := tracks.String("sortAttr", "", "sort results <attr>")
	tracksPerPage := tracks.Int("pp", 20, "per page")
	tracksPage := tracks.Int("p", 1, "page")
	tracksTag := tracks.String("tag", "", "filter by <tag>")
	tracksOutput := tracks.String("o", "text", "print output as text/json")

	path := args[0]
	switch path {
	case "eras":
		if err := eras.Parse(args[1:]); err != nil {
			return fmt.Errorf("error parsing eras args: %w", err)
		}
		c.PrintJSON = *erasOutput == "json"
	case "years":
		if err := years.Parse(args[1:]); err != nil {
            return err
        }
        // let's always include this
        c.Parameters = append(c.Parameters, "include_show_counts=true")
        c.PrintJSON = *yearsOutput == "json"
	case "songs":
		if err := songs.Parse(args[1:]); err != nil {
			return err
		}
		if err := c.validateOutput(*songsOutput, false); err != nil {
			return err
		}
		c.parseSortParams(*songsSortDir, *songsSortAttr)
		c.parsePageParams(*songsPerPage, *songsPage)
	case "tours":
		if err := tours.Parse(args[1:]); err != nil {
			return err
		}
		c.PrintJSON = *toursOutput == "json"
	case "venues":
		if err := venues.Parse(args[1:]); err != nil {
			return err
		}
		if err := c.validateOutput(*venuesOutput, false); err != nil {
			return err
		}
		c.parseSortParams(*venuesSortDir, *venuesSortAttr)
		c.parsePageParams(*venuesPerPage, *venuesPage)
	case "shows":
		if err := shows.Parse(args[1:]); err != nil {
			return err
		}
		if err := c.validateOutput(*showsOutput, *showsVerbose); err != nil {
			return err
		}
		c.parseSortParams(*showsSortDir, *showsSortAttr)
		c.parsePageParams(*showsPerPage, *showsPage)
		c.parseTag(*showsTag)
	case "show-on-date":
		if err := shows.Parse(args[1:]); err != nil {
			return err
		}
		if c.Query == "" {
			// todo put usage here
			return errors.New("need a date")
		}
		if err := c.validateOutput(*showsOutput, *showsVerbose); err != nil {
			return err
		}
	case "shows-on-day-of-year":
		if err := shows.Parse(args[1:]); err != nil {
			return err
		}
		if c.Query == "" {
			// todo put usage here
			return errors.New("need a day")
		}
		if err := c.validateOutput(*showsOutput, *showsVerbose); err != nil {
			return err
		}
	case "random-show":
		if err := shows.Parse(args[1:]); err != nil {
			return err
		}
		// doesn't take a parameter, so drop if user added one
		c.Query = ""
		if err := c.validateOutput(*showsOutput, *showsVerbose); err != nil {
			return err
		}
	case "tracks":
		if err := tracks.Parse(args[1:]); err != nil {
			return err
		}
		if err := c.validateOutput(*tracksOutput, false); err != nil {
			return err
		}
		c.parseSortParams(*tracksSortDir, *tracksSortAttr)
		c.parsePageParams(*tracksPerPage, *tracksPage)
		c.parseTag(*tracksTag)
	case "search":
		if c.Query == "" {
			// todo put usage here
			return errors.New("need a search term")
		}
	// case "playlists":
	case "tags":
		if err := tags.Parse(args[1:]); err != nil {
			return err
		}
		c.PrintJSON = *tagsOutput == "json"
	default:
		return fmt.Errorf("%s is not a recognized command", path)
	}
	return nil
}

func (c *Client) validateOutput(output string, verbose bool) error {
    // output
    if output != "text" && output != "json" {
        return errors.New("output must be text or json")
    }
    c.PrintJSON = output == "json"
	// verbose printing
	c.Verbose = verbose
    return nil
}

func (c *Client) parseSortParams(sortDir, sortAttr string) {
	switch sortDir {
    case "asc":
        c.Parameters = append(c.Parameters, "sort_dir=asc")
    case "desc":
        c.Parameters = append(c.Parameters, "sort_dir=desc")
    default:
        // just ignore
    }
    if sortAttr != "" {
        c.Parameters = append(c.Parameters, fmt.Sprintf("sort_attr=%s", sortAttr))
    }
}

func (c *Client) parsePageParams(perPage, page int) {
    if perPage != 20 && perPage > 0 {
        c.Parameters = append(c.Parameters, fmt.Sprintf("per_page=%d", perPage))
    }
    if page > 1 {
        c.Parameters = append(c.Parameters, fmt.Sprintf("page=%d", page))
    }
}

// todo match against possiblilities or just accept input?
func (c *Client) parseTag(tag string) {
	if tag != "" {
        c.Parameters = append(c.Parameters, fmt.Sprintf("tags=%s", tag))
    }
}

func (c *Client) Get(ctx context.Context, url string, data any) error {
	// url = "https://phish.in/api/v1/tracks?per_page=3&page=4&sort_dir=asc"
	// fmt.Println("url", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error building request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	authToken := c.APIKey
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	// todo add repo as ua header
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %q", resp.Status)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	// fmt.Println(string(b))
	// os.Exit(0)

	err = json.Unmarshal(b, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error decoding json response: %v\n", string(b))
		return err
	}
	return nil
}

func (c *Client) getAndPrintEras(ctx context.Context, url string) error {
	eras, err := c.getEras(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get eras data: %w", err)
	}
	if c.PrintJSON {
		return printJSON(c.Output, eras)
	}
	return prettyPrintEras(c.Output, eras)
}

func (c *Client) getEras(ctx context.Context, url string) (ErasOutput, error) {
	var resp ErasResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return ErasOutput{}, fmt.Errorf("unable to get eras list: %w", err)
	}
	o := ErasOutput{
		One: resp.Data.One,
		Two: resp.Data.Two,
		Three: resp.Data.Three,
		Four: resp.Data.Four,
	}
	return o, nil
}

func (c *Client) getAndPrintEra(ctx context.Context, url string) error {
	era, err := c.getEra(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get era data: %w", err)
	}
	if c.PrintJSON {
		return printJSON(c.Output, era)
	}
	return prettyPrintEra(c.Output, era)
}

func (c *Client) getEra(ctx context.Context, url string) (EraOutput, error) {
	var resp EraResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return EraOutput{}, fmt.Errorf("unable to get era details: %w", err)
	}
	o := EraOutput{
		EraName: c.Query,
		Years: resp.Era,
	}
	return o, nil
}

func (c *Client) getAndPrintYears(ctx context.Context, url string) error {
	years, err := c.getYears(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get years data: %w", err)
	}
	if c.PrintJSON {
		return printJSON(c.Output, years)
	}
	return prettyPrintYears(c.Tabwriter, years)
}

func (c *Client) getYears(ctx context.Context, url string) (YearsOutput, error) {
	var resp YearsResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return YearsOutput{}, fmt.Errorf("unable to get years list: %w", err)
	}
	o := YearsOutput{
		Years: resp.Data,
	}
	return o, nil
}

func (c *Client) getAndPrintYear(ctx context.Context, url string) error {
	shows, err := c.getYear(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get year data: %w", err)
	}
	if c.PrintJSON {
		return printJSON(c.Output, shows)
	}
	return prettyPrintShows(c.Tabwriter, shows, c.Verbose)
}

func (c *Client) getYear(ctx context.Context, url string) (ShowsOutput, error) {
	var resp YearResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return ShowsOutput{}, fmt.Errorf("unable to get year details: %w", err)
	}
	return convertToShowsOutput(resp.Data), nil
}

func (c *Client) getAndPrintShows(ctx context.Context, url string) error {
	shows, err := c.getShows(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get shows data: %w", err)
	}
	if c.PrintJSON {
		return printJSON(c.Output, shows)
	}
	return prettyPrintShows(c.Tabwriter, shows, c.Verbose)
}

func (c *Client) getShows(ctx context.Context, url string) (ShowsOutput, error) {
	var resp ShowsResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return ShowsOutput{}, fmt.Errorf("unable to get shows list: %w", err)
	}
	o := convertToShowsOutput(resp.Data)
	o.TotalEntries = resp.TotalEntries
	o.TotalPages = resp.TotalPages
	o.CurrentPage = resp.Page
	return o, nil
}

func (c *Client) getAndPrintShow(ctx context.Context, url string) error {
	show, err := c.getShow(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get show data: %w", err)
	}
	if c.PrintJSON {
		return printJSON(c.Output, show)
	}
	return prettyPrintShow(c.Tabwriter, show, c.Verbose)
}

func (c *Client) getShow(ctx context.Context, url string) (ShowOutput, error) {
	var resp ShowResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return ShowOutput{}, fmt.Errorf("unable to get show details: %w", err)
	}
	return convertToShowOutput(resp.Data), nil
}

func (c *Client) getAndPrintTours(ctx context.Context, url string) error {
	tours, err := c.getTours(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get tours data: %w", err)
	}
	if c.PrintJSON {
		return printJSON(c.Output, tours)
	}
	return prettyPrintTours(c.Tabwriter, tours)
}

func (c *Client) getTours(ctx context.Context, url string) (ToursOutput, error) {
	var resp ToursResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return ToursOutput{}, fmt.Errorf("unable to get tours list: %w", err)
	}
	tours := make([]TourOutput, 0, len(resp.Data))
	for _, t := range resp.Data {
		tour := TourOutput{
			Name: t.Name,
			ShowsCount: t.ShowsCount,
			StartsOn: t.StartsOn,
			EndsOn: t.EndsOn,
		}
		shows := convertToShowsOutput(t.Shows)
		tour.Shows = shows.Shows
		tours = append(tours, tour)
	}
	return ToursOutput{
		Tours: tours,
	}, nil
}

func (c *Client) getAndPrintTour(ctx context.Context, url string) error {
	tour, err := c.getTour(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get tour data: %w", err)
	}
	if c.PrintJSON {
		return printJSON(c.Output, tour)
	}
	return prettyPrintTour(c.Tabwriter, tour)
}

func (c *Client) getTour(ctx context.Context, url string) (TourOutput, error) {
	var resp TourResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return TourOutput{}, fmt.Errorf("unable to get tour details: %w", err)
	}
	o := TourOutput{
		Name: resp.Data.Name,
		ShowsCount: resp.Data.ShowsCount,
		StartsOn: resp.Data.StartsOn,
		EndsOn: resp.Data.EndsOn,
	}
	shows := convertToShowsOutput(resp.Data.Shows)
	o.Shows = shows.Shows
	return o, nil
}

func (c *Client) getAndPrintVenues(ctx context.Context, url string) error {
	venues, err := c.getVenues(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get venues data: %w", err)
	}
	if c.PrintJSON {
		return printJSON(c.Output, venues)
	}
	return prettyPrintVenues(c.Tabwriter, venues)
}

func (c *Client) getVenues(ctx context.Context, url string) (VenuesOutput, error) {
	var resp VenuesResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return VenuesOutput{}, fmt.Errorf("unable to get tours list: %w", err)
	}
	venues := make([]VenueOutput, 0, len(resp.Data))
	for _, v := range resp.Data {
		venues = append(venues, convertToVenueOutput(v))
	}
	return VenuesOutput{
		TotalEntries: resp.TotalEntries,
		TotalPages: resp.TotalPages,
		CurrentPage: resp.Page,
		Venues: venues,
	}, nil
}

func (c *Client) getAndPrintVenue(ctx context.Context, url string) error {
	venue, err := c.getVenue(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get venue data: %w", err)
	}
	if c.PrintJSON {
		return printJSON(c.Output, venue)
	}
	return prettyPrintVenue(c.Tabwriter, venue)
}

func (c *Client) getVenue(ctx context.Context, url string) (VenueOutput, error) {
	var resp VenueResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return VenueOutput{}, fmt.Errorf("unable to get tour details: %w", err)
	}
	return convertToVenueOutput(resp.Data), nil
}

func (c *Client) getAndPrintTags(ctx context.Context, url string) error {
	tags, err := c.getTags(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get tags data: %w", err)
	}
	if c.PrintJSON {
		return printJSON(c.Output, tags)
	}
	return prettyPrintTags(c.Tabwriter, tags)
}

func (c *Client) getTags(ctx context.Context, url string) (TagsOutput, error) {
	var resp TagsResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return TagsOutput{}, fmt.Errorf("unable to get tags list: %w", err)
	}
	tags := make([]TagListItemOutput, 0, len(resp.Data))
	for _, t := range resp.Data {
		tags = append(tags, TagListItemOutput{
			Name: t.Name,
			Group: t.Group,
			Description: t.Description,
			ShowIds: t.ShowIds,
			TrackIds: t.TrackIds,
		})
	}
	return TagsOutput{
		Tags: tags,
	}, nil
}

func (c *Client) getAndPrintTag(ctx context.Context, url string) error {
	tag, err := c.getTag(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get tag data: %w", err)
	}
	if c.PrintJSON {
		return printJSON(c.Output, tag)
	}
	return prettyPrintTag(c.Tabwriter, tag)
}

func (c *Client) getTag(ctx context.Context, url string) (TagListItemOutput, error) {
	var resp TagResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return TagListItemOutput{}, fmt.Errorf("unable to get tour details: %w", err)
	}
	o := TagListItemOutput{
		Name: resp.Data.Name,
		Group: resp.Data.Group,
		Description: resp.Data.Description,
		ShowIds: resp.Data.ShowIds,
		TrackIds: resp.Data.TrackIds,
	}
	return o, nil
}

func (c *Client) getAndPrintSongs(ctx context.Context, url string) error {
	songs, err := c.getSongs(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get songs data: %w", err)
	}
	if c.PrintJSON {
		return printJSON(c.Output, songs)
	}
	return prettyPrintSongs(c.Tabwriter, songs)
}

func (c *Client) getSongs(ctx context.Context, url string) (SongsOutput, error) {
	var resp SongsResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return SongsOutput{}, fmt.Errorf("unable to get songs list: %w", err)
	}
	songs := make([]SongOutput, 0, len(resp.Data))
	for _, s := range resp.Data {
		song := convertToSongOutput(s)
		songs = append(songs, song)
	}
	o := SongsOutput{
		TotalEntries: resp.TotalEntries,
		TotalPages: resp.TotalPages,
		CurrentPage: resp.Page,
		Songs: songs,
	}
	return o, nil
}

func (c *Client) getAndPrintSong(ctx context.Context, url string) error {
	song, err := c.getSong(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get song data: %w", err)
	}
	if c.PrintJSON {
		return printJSON(c.Output, song)
	}
	return prettyPrintSong(c.Tabwriter, song)
}

func (c *Client) getSong(ctx context.Context, url string) (SongOutput, error) {
	var resp SongResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return SongOutput{}, fmt.Errorf("unable to get song details: %w", err)
	}
	return convertToSongOutput(resp.Data), nil
}

func (c *Client) getAndPrintTracks(ctx context.Context, url string) error {
	tracks, err := c.getTracks(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get tracks data: %w", err)
	}
	if c.PrintJSON {
		return printJSON(c.Output, tracks)
	}
	return prettyPrintTracks(c.Tabwriter, tracks)
}

func (c *Client) getTracks(ctx context.Context, url string) (TracksOutput, error) {
	var resp TracksResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return TracksOutput{}, fmt.Errorf("unable to get tracks list: %w", err)
	}
	o := convertToTracksOutput(resp.Data)
	o.TotalEntries = resp.TotalEntries
	o.TotalPages = resp.TotalPages
	o.CurrentPage = resp.Page
	return o, nil
}

func (c *Client) getAndPrintTrack(ctx context.Context, url string) error {
	track, err := c.getTrack(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get track data: %w", err)
	}
	if c.PrintJSON {
		return printJSON(c.Output, track)
	}
	return prettyPrintTrack(c.Tabwriter, track)
}

func (c *Client) getTrack(ctx context.Context, url string) (TrackOutput, error) {
	var resp TrackResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return TrackOutput{}, fmt.Errorf("unable to get track details: %w", err)
	}
	return convertToTrackOutput(resp.Data), nil
}
