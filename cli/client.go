package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"golang.org/x/sync/errgroup"
)

type Client struct {
	Tabwriter  *tabwriter.Writer
	HTTPClient *http.Client
	ErrGroup *errgroup.Group
	BaseURL    string
	APIKey     string
	PrintJSON  bool
	Query      string
	Parameters []string
	Output     io.Writer
	Verbose    bool
	Debug      bool
	Download bool
}

func NewClient(apiKey string, output io.Writer) *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
		BaseURL:    "https://phish.in/api/v1",
		APIKey:     apiKey,
		Output:     output,
		Tabwriter:  tabwriter.NewWriter(output, 0, 4, 2, ' ', tabwriter.DiscardEmptyColumns),
		ErrGroup: &errgroup.Group{},
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
	phishin := flag.NewFlagSet("phishin", flag.ExitOnError)
	query := phishin.String("search", "", "search query")
	phishin.StringVar(query, "s", "", "search query")
	output := phishin.String("output", "text", "print output as <text> or <json>")
	phishin.StringVar(output, "o", "text", "print output as <text> or <json>")
	sortDir := phishin.String("sort-dir", "", "sort results <asc> or <desc>")
	phishin.StringVar(sortDir, "dir", "", "sort results <asc> or <desc>")
	sortAttr := phishin.String("sort-attr", "", "sort results <attr>")
	phishin.StringVar(sortAttr, "a", "", "sort results <attr>")
	perPage := phishin.Int("per-page", 20, "number of results included per page")
	phishin.IntVar(perPage, "pp", 20, "number of results included per page")
	page := phishin.Int("page", 1, "result page to return")
	phishin.IntVar(page, "p", 1, "result page to return")
	tag := phishin.String("tag", "", "filter by <tag>")
	phishin.StringVar(tag, "t", "", "filter by <tag>")
	verbose := phishin.Bool("verbose", false, "verbose output")
	phishin.BoolVar(verbose, "v", false, "verbose output")
	debug := phishin.Bool("debug", false, "print the url that the client is sending to the server")
	download := phishin.Bool("d", false, "download (if applicable)")

	phishin.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		fmt.Println("Flags:")
		phishin.PrintDefaults()
	}
	if err := phishin.Parse(args[1:]); err != nil {
		return fmt.Errorf("error parsing args: %w", err)
	}

	c.Query = *query
	c.PrintJSON = *output == "json"
	c.Verbose = *verbose
	c.Debug = *debug
	c.Download = *download

	path := args[0]
	switch path {
	case showsPath, tracksPath:
		c.parseTag(*tag)
		c.parsePageParams(*perPage, *page)
		c.parseSortParams(*sortDir, *sortAttr)
	case songsPath, venuesPath:
		c.parsePageParams(*perPage, *page)
		c.parseSortParams(*sortDir, *sortAttr)
	case yearsPath:
		// let's always include this
		c.Parameters = append(c.Parameters, "include_show_counts=true")
	case showOnDatePath:
		if c.Query == "" {
			return errors.New("need a date")
		}
	case showsDayOfYearPath:
		if c.Query == "" {
			return errors.New("need a day")
		}
	case randomShowPath:
		// doesn't take a parameter, so drop if user added one
		c.Query = ""
	case searchPath:
		if c.Query == "" {
			return errors.New("need a search term")
		}
	case erasPath, toursPath, tagsPath:
		// do nothing
	default:
		fmt.Fprintf(os.Stderr, "%s is not a recognized command\n", path)
		return errors.New(endpointList)
	}
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

func (c *Client) parseTag(tag string) {
	if tag != "" {
		c.Parameters = append(c.Parameters, fmt.Sprintf("tag=%s", tag))
	}
}

func (c *Client) Get(ctx context.Context, url string, data any) error {
	if c.Debug {
		fmt.Fprintln(c.Output, url)
	}
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

	err = json.Unmarshal(b, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error decoding json response: %v\n", string(b))
		return err
	}
	return nil
}

func (c *Client) run(ctx context.Context, path string) error {
	url := c.FormatURL(path)
	switch path {
	case erasPath:
		if c.Query != "" {
			if err := c.getAndPrintEra(ctx, url); err != nil {
				return fmt.Errorf("era details failure: %w", err)
			}
			return nil
		}
		if err := c.getAndPrintEras(ctx, url); err != nil {
			return fmt.Errorf("eras list failure: %w", err)
		}
		return nil
	case yearsPath:
		if c.Query != "" {
			if err := c.getAndPrintYear(ctx, url); err != nil {
				return fmt.Errorf("year details failure: %w", err)
			}
			return nil
		}
		if err := c.getAndPrintYears(ctx, url); err != nil {
			return fmt.Errorf("years list failure: %w", err)
		}
		return nil
	case songsPath:
		if c.Query != "" {
			if err := c.getAndPrintSong(ctx, url); err != nil {
				return fmt.Errorf("song details failure: %w", err)
			}
			return nil
		}
		if err := c.getAndPrintSongs(ctx, url); err != nil {
			return fmt.Errorf("songs list failure: %w", err)
		}
		return nil
	case toursPath:
		if c.Query != "" {
			if err := c.getAndPrintTour(ctx, url); err != nil {
				return fmt.Errorf("tour details failure: %w", err)
			}
			return nil
		}
		if err := c.getAndPrintTours(ctx, url); err != nil {
			return fmt.Errorf("tours list failure: %w", err)
		}
		return nil
	case venuesPath:
		if c.Query != "" {
			if err := c.getAndPrintVenue(ctx, url); err != nil {
				return fmt.Errorf("venue details failure: %w", err)
			}
			return nil
		}
		if err := c.getAndPrintVenues(ctx, url); err != nil {
			return fmt.Errorf("venues list failure: %w", err)
		}
		return nil
	case showsPath:
		if c.Query != "" {
			if err := c.getAndPrintShow(ctx, url); err != nil {
				return fmt.Errorf("show details failure: %w", err)
			}
			return nil
		}
		if err := c.getAndPrintShows(ctx, url); err != nil {
			return fmt.Errorf("shows list failure: %w", err)
		}
		return nil
	case showOnDatePath:
		if err := c.getAndPrintShow(ctx, url); err != nil {
			return fmt.Errorf("show details failure: %w", err)
		}
		return nil
	case showsDayOfYearPath:
		if err := c.getAndPrintShows(ctx, url); err != nil {
			return fmt.Errorf("shows list failure: %w", err)
		}
		return nil
	case randomShowPath:
		if err := c.getAndPrintShow(ctx, url); err != nil {
			return fmt.Errorf("show details failure: %w", err)
		}
		return nil
	case tracksPath:
		if c.Query != "" {
			if err := c.getAndPrintTrack(ctx, url); err != nil {
				return fmt.Errorf("track details failure: %w", err)
			}
			return nil
		}
		if err := c.getAndPrintTracks(ctx, url); err != nil {
			return fmt.Errorf("tracks list failure: %w", err)
		}
		return nil
	case searchPath:
		// todo
		if err := c.Get(ctx, url, nil); err != nil {
			return err
		}
		return nil
	// case "playlists":

	case tagsPath:
		if c.Query != "" {
			if err := c.getAndPrintTag(ctx, url); err != nil {
				return fmt.Errorf("tag details failure: %w", err)
			}
			return nil
		}
		if err := c.getAndPrintTags(ctx, url); err != nil {
			return fmt.Errorf("tags list failure: %w", err)
		}
		return nil
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
		One:   resp.Data.One,
		Two:   resp.Data.Two,
		Three: resp.Data.Three,
		Four:  resp.Data.Four,
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
		Years:   resp.Era,
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
	return convertShowToShowsOutput(resp.Data), nil
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
	o := convertShowToShowsOutput(resp.Data)
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
	if c.Download {
		for _, t := range resp.Data.Tracks {
			c.DownloadTrack(ctx, t.Mp3, t.Slug)
		}
	}
	
	return convertShowToShowOutput(resp.Data), nil
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
			Name:       t.Name,
			ShowsCount: t.ShowsCount,
			StartsOn:   t.StartsOn,
			EndsOn:     t.EndsOn,
		}
		shows := convertShowToShowsOutput(t.Shows)
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
		Name:       resp.Data.Name,
		ShowsCount: resp.Data.ShowsCount,
		StartsOn:   resp.Data.StartsOn,
		EndsOn:     resp.Data.EndsOn,
	}
	shows := convertShowToShowsOutput(resp.Data.Shows)
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
		venues = append(venues, convertVenueToVenueOutput(v))
	}
	return VenuesOutput{
		TotalEntries: resp.TotalEntries,
		TotalPages:   resp.TotalPages,
		CurrentPage:  resp.Page,
		Venues:       venues,
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
		return VenueOutput{}, fmt.Errorf("unable to get venue details: %w", err)
	}
	return convertVenueToVenueOutput(resp.Data), nil
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
			Name:        t.Name,
			Group:       t.Group,
			Description: t.Description,
			ShowIds:     t.ShowIds,
			TrackIds:    t.TrackIds,
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
		Name:        resp.Data.Name,
		Group:       resp.Data.Group,
		Description: resp.Data.Description,
		ShowIds:     resp.Data.ShowIds,
		TrackIds:    resp.Data.TrackIds,
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
		song := convertSongToSongOutput(s)
		songs = append(songs, song)
	}
	o := SongsOutput{
		TotalEntries: resp.TotalEntries,
		TotalPages:   resp.TotalPages,
		CurrentPage:  resp.Page,
		Songs:        songs,
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
	return convertSongToSongOutput(resp.Data), nil
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
	o := convertTracksToTracksOutput(resp.Data)
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
	if c.Download {
		c.ErrGroup.Go(func() error {
			fileName := fmt.Sprintf("%s.mp3", resp.Data.Slug)
			return c.DownloadTrack(ctx, resp.Data.Mp3, fileName)
		})
	}
	return convertTrackToTrackOutput(resp.Data), nil
}

// todo think about cleanup for cancelled context?
func (c *Client) DownloadTrack(ctx context.Context, url, fileName string) error {
	time.Sleep(4 * time.Second)
	p := filepath.Join(c.Query, fileName)
	f, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() { _ =  f.Close()}()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get response: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received unexpected status code: %q", resp.Status)
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("unable to copy data to file: %w", err)
	}
	return nil
}