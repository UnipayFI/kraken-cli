package exchange

import (
	kraken "github.com/UnipayFI/go-kraken"
	"github.com/UnipayFI/go-kraken/client"
	"github.com/UnipayFI/kraken-cli/config"
)

// Client wraps the go-kraken REST client construction, applying credentials and
// optional base URL / proxy from the process config.
type Client struct {
	APIKey    string
	APISecret string
	BaseURL   string
	Proxy     string
}

func NewClient(apiKey, apiSecret string) *Client {
	return &Client{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   config.Config.BaseURL,
		Proxy:     config.Config.Proxy,
	}
}

func (c *Client) options() []client.Options {
	opts := []client.Options{client.WithAuth(c.APIKey, c.APISecret)}
	if c.BaseURL != "" {
		opts = append(opts, client.WithBaseURL(c.BaseURL))
	}
	if c.Proxy != "" {
		opts = append(opts, client.WithProxy(c.Proxy))
	}
	return opts
}

// NewKrakenClient builds a fresh go-kraken REST client for a single command
// invocation.
func (c *Client) NewKrakenClient() *kraken.Client {
	return kraken.NewClient(c.options()...)
}
