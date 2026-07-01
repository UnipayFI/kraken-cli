package cmd

import (
	"github.com/UnipayFI/kraken-cli/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version and build information",
	Run: func(cmd *cobra.Command, args []string) {
		version.Version()
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
