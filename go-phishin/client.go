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
	eras := flag.NewFlagSet("eras", flag.ExitOnError)
	eras.StringVar(&c.Query, "s", "", "search query")
	erasOutput := eras.String("o", "text", "print output as text/json")

	years := flag.NewFlagSet("years", flag.ExitOnError)
	years.StringVar(&c.Query, "s", "", "search query")
	yearsSortDir := years.String("sortDir", "", "sort results asc/desc") 
    yearsSortAttr := years.String("sortAttr", "", "sort results <attr>")
	yearsPerPage := years.Int("pp", 20, "per page") // todo check
	yearsPage := years.Int("p", 1, "page")
	yearsVerbose := years.Bool("v", false, "print setlists (valid only when searching for a specific year)")
    yearsOutput := years.String("o", "text", "print output as text/json")
    
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
        if err := c.validateParams(*yearsOutput, *yearsVerbose, *yearsSortDir, *yearsSortAttr, "", *yearsPerPage, *yearsPage); err != nil {
            return err
        }
        // let's always include this
        c.Parameters = append(c.Parameters, "include_show_counts=true")
        c.PrintJSON = *yearsOutput == "json"
	case "songs":
	case "tours":
	case "venues":
	case "shows":
	case "show-on-date":
	case "shows-on-day-of-year":
	case "random-show":
	case "tracks":
	case "search":
	// case "playlists":
	case "tags":
	default:
		return fmt.Errorf("%s is not a recognized command", path)
	}
	return nil
}

func (c *Client) validateParams(output string, verbose bool, sortDir, sortAttr, tag string, perPage, page int) error {
    // output
    if output != "text" && output != "json" {
        return errors.New("output must be text or json")
    }
    c.PrintJSON = output == "json"
	// verbose printing
	c.Verbose = verbose
    // sortDir
    switch sortDir {
    case "asc":
        c.Parameters = append(c.Parameters, "sort_dir=asc")
    case "desc":
        c.Parameters = append(c.Parameters, "sort_dir=desc")
    default:
        // just ignore
    }
    // sortAttr
    if sortAttr != "" {
        c.Parameters = append(c.Parameters, fmt.Sprintf("sort_attr=%s", sortAttr))
    }
    if tag != "" {
        c.Parameters = append(c.Parameters, fmt.Sprintf("tags=%s", tag))
    }
    // perPage
    if perPage != 20 && perPage > 0 {
        c.Parameters = append(c.Parameters, fmt.Sprintf("per_page=%d", perPage))
    }
    // page
    if page > 1 {
        c.Parameters = append(c.Parameters, fmt.Sprintf("page=%d", page))
    }
    return nil
}

func (c *Client) Get(ctx context.Context, url string, data any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error building request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	authToken := c.APIKey
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
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

	// todo remove
	// fmt.Println(string(b))
	// os.WriteFile("foo", b, 0644)
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
		return printJSONEras(c.Output, eras)
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
		return printJSONEra(c.Output, era)
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
		return printJSONYears(c.Output, years)
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

func (c *Client) getYear(ctx context.Context, url string) (ShowsOutput, error) {
	var resp YearResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return ShowsOutput{}, fmt.Errorf("unable to get year details: %w", err)
	}
	o := ShowsOutput{
		Shows: resp.Data,
	}
	return o, nil
}

func (c *Client) getAndPrintYear(ctx context.Context, url string) error {
	shows, err := c.getYear(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get year data: %w", err)
	}
	if c.PrintJSON {
		return printJSONShows(c.Output, shows)
	}
	return prettyPrintShows(c.Tabwriter, shows, c.Verbose)
}

func (c *Client) getAndPrintShows(ctx context.Context, url string) error {
	shows, err := c.getShows(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get shows data: %w", err)
	}
	if c.PrintJSON {
		return printJSONShows(c.Output, shows)
	}
	return prettyPrintShows(c.Tabwriter, shows, c.Verbose)
}

func (c *Client) getShows(ctx context.Context, url string) (ShowsOutput, error) {
	var resp ShowsResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return ShowsOutput{}, fmt.Errorf("unable to get shows list: %w", err)
	}

	o := ShowsOutput{
		Shows: resp.Data,
	}
	return o, nil
}