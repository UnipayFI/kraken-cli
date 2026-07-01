package version

import "fmt"

// Set via -ldflags at build time.
var (
	version   string
	buildTime string
)

func Version() {
	fmt.Printf("Version: %s, Build Time: %s\n", version, buildTime)
	fmt.Println("Author: feeeei unipay")
}

// String returns the version string for cobra's --version flag.
func String() string {
	if version == "" {
		return "dev"
	}
	return version
}
