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

	"golang.org/x/sync/errgroup"
)

type Client struct {
	HTTPClient *http.Client
	ErrGroup   *errgroup.Group
	BaseURL    string
	APIKey     string
	PrintJSON  bool
	Query      string
	Parameters []string
	Output     io.Writer
	Verbose    bool
	Debug      bool
	Download   bool
	RawOutput  bool
}

func NewClient(apiKey string, output io.Writer) *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
		BaseURL:    "https://phish.in/api/v1",
		APIKey:     apiKey,
		Output:     output,
		ErrGroup:   &errgroup.Group{},
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
	raw := phishin.Bool("raw", false, "print full api json response")
	phishin.BoolVar(raw, "r", false, "print full api json response")

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
	c.RawOutput = *raw

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

func (c *Client) getAndPrintRaw(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error building request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	authToken := c.APIKey
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	req.Header.Set("User-Agent", "https://github.com/davemolk/phishin")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			fmt.Fprint(os.Stderr, searchTips)
		}
		return fmt.Errorf("unexpected response status: %q", resp.Status)
	}
	g := &GenericResponse{}
	json.NewDecoder(resp.Body).Decode(g)
	b, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal json: %w", err)
	}
	fmt.Fprintln(tabwriter.NewWriter(c.Output, 0, 4, 2, ' ', tabwriter.DiscardEmptyColumns), string(b))
	return nil
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
	req.Header.Set("User-Agent", "https://github.com/davemolk/phishin")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			fmt.Fprint(os.Stderr, searchTips)
		}
		return fmt.Errorf("unexpected response status: %q", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(data)
}

func (c *Client) run(ctx context.Context, path string) error {
	url := c.FormatURL(path)
	if c.RawOutput {
		return c.getAndPrintRaw(ctx, url)
	}
	var results PrettyPrinter
	var err error
	switch {
	case path == erasPath && c.Query != "":
		results, err = c.getEra(ctx, url)
		if err != nil {
			return fmt.Errorf("era details failure: %w", err)
		}
	case path == erasPath:
		results, err = c.getEras(ctx, url)
		if err != nil {
			return fmt.Errorf("eras list failure: %w", err)
		}
	case path == yearsPath && c.Query != "":
		results, err = c.getYear(ctx, url)
		if err != nil {
			return fmt.Errorf("year details failure: %w", err)
		}
	case path == yearsPath:
		results, err = c.getYears(ctx, url)
		if err != nil {
			return fmt.Errorf("years list failure: %w", err)
		}
	case path == songsPath && c.Query != "":
		results, err = c.getSong(ctx, url)
		if err != nil {
			return fmt.Errorf("song details failure: %w", err)
		}
	case path == songsPath:
		results, err = c.getSongs(ctx, url)
		if err != nil {
			return fmt.Errorf("songs list failure: %w", err)
		}
	case path == toursPath && c.Query != "":
		results, err = c.getTour(ctx, url)
		if err != nil {
			return fmt.Errorf("tour details failure: %w", err)
		}
	case path == toursPath:
		results, err = c.getTours(ctx, url)
		if err != nil {
			return fmt.Errorf("tours list failure: %w", err)
		}
	case path == venuesPath && c.Query != "":
		results, err = c.getVenue(ctx, url)
		if err != nil {
			return fmt.Errorf("venue details failure: %w", err)
		}
	case path == venuesPath:
		results, err = c.getVenues(ctx, url)
		if err != nil {
			return fmt.Errorf("venues list failure: %w", err)
		}
	case path == showsPath && c.Query != "":
		results, err = c.getShow(ctx, url)
		if err != nil {
			return fmt.Errorf("show details failure: %w", err)
		}
	// todo consolidate these
	case path == showsPath:
		results, err = c.getShows(ctx, url)
		if err != nil {
			return fmt.Errorf("shows list failure: %w", err)
		}
	case path == showOnDatePath:
		results, err = c.getShow(ctx, url)
		if err != nil {
			return fmt.Errorf("show details failure: %w", err)
		}
	case path == showsDayOfYearPath:
		results, err = c.getShows(ctx, url)
		if err != nil {
			return fmt.Errorf("shows list failure: %w", err)
		}
	case path == randomShowPath:
		results, err = c.getShow(ctx, url)
		if err != nil {
			return fmt.Errorf("show details failure: %w", err)
		}
	case path == tracksPath && c.Query != "":
		results, err = c.getTrack(ctx, url)
		if err != nil {
			return fmt.Errorf("track details failure: %w", err)
		}
	case path == tracksPath:
		results, err = c.getTracks(ctx, url)
		if err != nil {
			return fmt.Errorf("tracks list failure: %w", err)
		}
	case path == searchPath:
		results, err = c.getSearch(ctx, url)
		if err != nil {
			return fmt.Errorf("search failure: %w", err)
		}
	// case path == "playlists" && c.Query != "":

	case path == tagsPath && c.Query != "":
		results, err = c.getTag(ctx, url)
		if err != nil {
			return fmt.Errorf("tag details failure: %w", err)
		}
	case path == tagsPath:
		results, err = c.getTags(ctx, url)
		if err != nil {
			return fmt.Errorf("tags list failure: %w", err)
		}
	}
	return PrintResults(c.Output, results, c.PrintJSON, c.Verbose)
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

func (c *Client) getYear(ctx context.Context, url string) (ShowsOutput, error) {
	var resp YearResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return ShowsOutput{}, fmt.Errorf("unable to get year details: %w", err)
	}
	return convertShowsToOutput(resp.Data), nil
}

func (c *Client) getShows(ctx context.Context, url string) (ShowsOutput, error) {
	var resp ShowsResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return ShowsOutput{}, fmt.Errorf("unable to get shows list: %w", err)
	}
	o := convertShowsToOutput(resp.Data)
	o.TotalEntries = resp.TotalEntries
	o.TotalPages = resp.TotalPages
	o.CurrentPage = resp.Page
	return o, nil
}

func (c *Client) getShow(ctx context.Context, url string) (ShowOutput, error) {
	var resp ShowResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return ShowOutput{}, fmt.Errorf("unable to get show details: %w", err)
	}
	if c.Download {
		if err := os.Mkdir(resp.Data.Date, 0755); err != nil {
			return ShowOutput{}, fmt.Errorf("unable to create directory for downloaded files: %w", err)
		}
		for i, t := range resp.Data.Tracks {
			// start track number with 1, capture loop vars locally
			i, t := i+1, t
			c.ErrGroup.Go(func() error {
				fileName := fmt.Sprintf("%d-%s.mp3", i, t.Slug)
				return c.DownloadTrack(ctx, t.Mp3, fileName, resp.Data.Date)
			})
		}
	}
	return convertShowToOutput(resp.Data), nil
}

func (c *Client) getTours(ctx context.Context, url string) (ToursOutput, error) {
	var resp ToursResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return ToursOutput{}, fmt.Errorf("unable to get tours list: %w", err)
	}
	return convertToursToOutput(resp.Data), nil
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
	shows := convertShowsToOutput(resp.Data.Shows)
	o.Shows = shows.Shows
	return o, nil
}

func (c *Client) getVenues(ctx context.Context, url string) (VenuesOutput, error) {
	var resp VenuesResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return VenuesOutput{}, fmt.Errorf("unable to get tours list: %w", err)
	}
	venues := make([]VenueOutput, 0, len(resp.Data))
	for _, v := range resp.Data {
		venues = append(venues, convertVenueToOutput(v))
	}
	return VenuesOutput{
		TotalEntries: resp.TotalEntries,
		TotalPages:   resp.TotalPages,
		CurrentPage:  resp.Page,
		Venues:       venues,
	}, nil
}

func (c *Client) getVenue(ctx context.Context, url string) (VenueOutput, error) {
	var resp VenueResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return VenueOutput{}, fmt.Errorf("unable to get venue details: %w", err)
	}
	return convertVenueToOutput(resp.Data), nil
}

func (c *Client) getTags(ctx context.Context, url string) (TagsOutput, error) {
	var resp TagsResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return TagsOutput{}, fmt.Errorf("unable to get tags list: %w", err)
	}
	tags := make([]TagListItemOutput, 0, len(resp.Data))
	for _, t := range resp.Data {
		tags = append(tags, convertTagListItemToOutput(t))
	}
	return TagsOutput{
		Tags: tags,
	}, nil
}

func (c *Client) getTag(ctx context.Context, url string) (TagListItemOutput, error) {
	var resp TagResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return TagListItemOutput{}, fmt.Errorf("unable to get tour details: %w", err)
	}
	return convertTagListItemToOutput(resp.Data), nil
}

func (c *Client) getSongs(ctx context.Context, url string) (SongsOutput, error) {
	var resp SongsResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return SongsOutput{}, fmt.Errorf("unable to get songs list: %w", err)
	}
	songs := make([]SongOutput, 0, len(resp.Data))
	for _, s := range resp.Data {
		song := convertSongToOutput(s)
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

func (c *Client) getSong(ctx context.Context, url string) (SongOutput, error) {
	var resp SongResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return SongOutput{}, fmt.Errorf("unable to get song details: %w", err)
	}
	return convertSongToOutput(resp.Data), nil
}

func (c *Client) getTracks(ctx context.Context, url string) (TracksOutput, error) {
	var resp TracksResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return TracksOutput{}, fmt.Errorf("unable to get tracks list: %w", err)
	}
	o := convertTracksToOutput(resp.Data)
	o.TotalEntries = resp.TotalEntries
	o.TotalPages = resp.TotalPages
	o.CurrentPage = resp.Page
	return o, nil
}

func (c *Client) getTrack(ctx context.Context, url string) (TrackOutput, error) {
	var resp TrackResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return TrackOutput{}, fmt.Errorf("unable to get track details: %w", err)
	}
	if c.Download {
		c.ErrGroup.Go(func() error {
			fileName := fmt.Sprintf("%s.mp3", resp.Data.Slug)
			return c.DownloadTrack(ctx, resp.Data.Mp3, fileName, ".")
		})
	}
	return convertTrackToOutput(resp.Data), nil
}

func (c *Client) getSearch(ctx context.Context, url string) (SearchOutput, error) {
	var resp SearchResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return SearchOutput{}, fmt.Errorf("couldn't get search results: %w", err)
	}
	return convertSearchToSearchOutput(resp), nil
}

// todo think about cleanup for cancelled context?
// todo handle progress counter differently when have concurrent downloads?
// todo track percentage via ContentLength
func (c *Client) DownloadTrack(ctx context.Context, url, fileName, dirName string) error {
	p := filepath.Join(dirName, fileName)
	f, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() { _ = f.Close() }()

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

	progress := &WriteCounter{
		Name: fileName,
	}
	_, err = io.Copy(f, io.TeeReader(resp.Body, progress))
	fmt.Println()
	if err != nil {
		return fmt.Errorf("unable to copy data to file: %w", err)
	}
	return nil
}
