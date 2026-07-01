package config

import "os"

// Config holds the process-wide CLI configuration. Credentials and connection
// settings come from the environment; OutputJSON is bound to the global --json
// flag in cmd/root.go.
var Config struct {
	APIKey     string
	APISecret  string
	BaseURL    string
	Proxy      string
	OutputJSON bool
}

func init() {
	Config.APIKey = os.Getenv("KRAKEN_API_KEY")
	Config.APISecret = os.Getenv("KRAKEN_API_SECRET")
	Config.BaseURL = os.Getenv("KRAKEN_BASE_URL")
	Config.Proxy = os.Getenv("KRAKEN_PROXY")
}
