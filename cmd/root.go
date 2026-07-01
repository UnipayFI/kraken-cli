package cmd

import (
	"errors"

	"github.com/UnipayFI/kraken-cli/config"
	"github.com/UnipayFI/kraken-cli/version"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:               "kraken-cli",
	Short:             "Kraken Spot API for CLI version",
	Version:           version.String(),
	PersistentPreRunE: checkCredentials,
}

func init() {
	initCommandConfig()
}

func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func initCommandConfig() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.PersistentFlags().BoolVar(&config.Config.OutputJSON, "json", false, "Output JSON instead of a table")
}

func checkCredentials(cmd *cobra.Command, args []string) error {
	if config.Config.APIKey == "" || config.Config.APISecret == "" {
		return errors.New("KRAKEN_API_KEY and KRAKEN_API_SECRET must be set")
	}
	return nil
}
