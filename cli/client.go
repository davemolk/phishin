package cli

import (
	"fmt"
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