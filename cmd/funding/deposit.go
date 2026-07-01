package funding

import (
	"log"

	"github.com/UnipayFI/kraken-cli/config"
	"github.com/UnipayFI/kraken-cli/exchange"
	"github.com/UnipayFI/kraken-cli/exchange/funding"
	"github.com/UnipayFI/kraken-cli/printer"
	"github.com/spf13/cobra"
)

func newClient() funding.Client {
	return funding.Client{Client: exchange.NewClient(config.Config.APIKey, config.Config.APISecret)}
}

var (
	depositMethodsCmd = &cobra.Command{
		Use:   "deposit-methods",
		Short: "List deposit methods for an asset",
		Long: `Retrieve the methods available for depositing a particular asset.
Required: --asset.

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-deposit-methods`,
		Run: showDepositMethods,
	}

	depositAddressesCmd = &cobra.Command{
		Use:   "deposit-addresses",
		Short: "List (or generate) deposit addresses",
		Long: `Retrieve (or generate, with --new) deposit addresses for a particular
asset and method. Required: --asset, --method.

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-deposit-addresses`,
		Run: showDepositAddresses,
	}

	depositStatusCmd = &cobra.Command{
		Use:   "deposit-status",
		Short: "Show status of recent deposits",
		Long: `Retrieve information about recent deposits. Results can be filtered by
--asset and --method.

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-status-recent-deposits`,
		Run: showDepositStatus,
	}
)

// InitDepositCmds returns the deposit-related subcommands.
func InitDepositCmds() []*cobra.Command {
	depositMethodsCmd.Flags().StringP("asset", "a", "", "asset, e.g. XBT (required)")
	depositMethodsCmd.MarkFlagRequired("asset")

	depositAddressesCmd.Flags().StringP("asset", "a", "", "asset, e.g. XBT (required)")
	depositAddressesCmd.Flags().StringP("method", "m", "", "deposit method, e.g. Bitcoin (required)")
	depositAddressesCmd.Flags().BoolP("new", "n", false, "generate a new address")
	depositAddressesCmd.MarkFlagRequired("asset")
	depositAddressesCmd.MarkFlagRequired("method")

	depositStatusCmd.Flags().StringP("asset", "a", "", "filter by asset")
	depositStatusCmd.Flags().StringP("method", "m", "", "filter by method")

	return []*cobra.Command{depositMethodsCmd, depositAddressesCmd, depositStatusCmd}
}

func showDepositMethods(cmd *cobra.Command, args []string) {
	client := newClient()
	asset, _ := cmd.Flags().GetString("asset")
	methods, err := client.GetDepositMethods(asset)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(methods)
}

func showDepositAddresses(cmd *cobra.Command, args []string) {
	client := newClient()
	asset, _ := cmd.Flags().GetString("asset")
	method, _ := cmd.Flags().GetString("method")
	newAddr, _ := cmd.Flags().GetBool("new")
	addrs, err := client.GetDepositAddresses(asset, method, newAddr)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(addrs)
}

func showDepositStatus(cmd *cobra.Command, args []string) {
	client := newClient()
	asset, _ := cmd.Flags().GetString("asset")
	method, _ := cmd.Flags().GetString("method")
	statuses, err := client.GetDepositStatus(asset, method)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(statuses)
}
