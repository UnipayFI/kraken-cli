package exchange

import (
	kraken "github.com/UnipayFI/go-kraken"
	"github.com/UnipayFI/go-kraken/client"
	"github.com/UnipayFI/kraken-cli/config"
)

// Client wraps the go-kraken REST client construction, applying credentials and
// an optional base URL override from the process config. Proxying follows the
// standard HTTPS_PROXY / HTTP_PROXY / NO_PROXY environment variables, which the
// underlying transport honours automatically.
type Client struct {
	APIKey    string
	APISecret string
	BaseURL   string
}

func NewClient(apiKey, apiSecret string) *Client {
	return &Client{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   config.Config.BaseURL,
	}
}

func (c *Client) options() []client.Options {
	opts := []client.Options{client.WithAuth(c.APIKey, c.APISecret)}
	if c.BaseURL != "" {
		opts = append(opts, client.WithBaseURL(c.BaseURL))
	}
	return opts
}

// NewKrakenClient builds a fresh go-kraken REST client for a single command
// invocation.
func (c *Client) NewKrakenClient() *kraken.Client {
	return kraken.NewClient(c.options()...)
}
