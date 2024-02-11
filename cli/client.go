package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL string
	APIKey string
	PrintJSON bool
	Query string
	Parameters []string
}

func NewClient(apiKey string) *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
		BaseURL: "https://phish.in/api/v1",
		APIKey: apiKey,
	}
}

func (c *Client) FormatURL(endpoint string) string {
	if c.Query != "" {
		// return now to avoid mixing in params
		return fmt.Sprintf("%s/%s/%s", c.BaseURL, endpoint, c.Query)
	}
	url := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)
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

	endpoint := args[0]
	switch endpoint {
	case "eras":
		if err := eras.Parse(args[1:]); err != nil {
			return fmt.Errorf("error parsing eras args: %w", err)
		}
		c.PrintJSON = *erasOutput == "json"
	default:
		return fmt.Errorf("%s is not a recognized command", endpoint)
	}
	return nil
}

func (c *Client) GetEras() (ErasOutput, error) {
	var resp ErasResponse
	if err := c.Get("eras", &resp); err != nil {
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

func (c *Client) GetEra() (EraOutput, error) {
	var resp EraResponse
	if err := c.Get("eras", &resp); err != nil {
		return EraOutput{}, fmt.Errorf("unable to get era details: %w", err)
	}
	o := EraOutput{
		Era: c.Query,
		EraList: resp.Era,
	}
	return o, nil
}

func (c *Client) Get(endpoint string, data any) error {
	ctx := context.Background()
	url := c.FormatURL(endpoint)
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
	if len(b) == 0 {
		return errors.New("no response body")
	}
	return json.Unmarshal(b, data)
}