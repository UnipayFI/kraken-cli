package funding

import (
	"log"

	"github.com/UnipayFI/kraken-cli/printer"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var (
	withdrawMethodsCmd = &cobra.Command{
		Use:   "withdraw-methods",
		Short: "List withdrawal methods",
		Long: `Retrieve the withdrawal methods available for the account, optionally
filtered by --asset and --network.

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-withdrawal-methods`,
		Run: showWithdrawMethods,
	}

	withdrawAddressesCmd = &cobra.Command{
		Use:   "withdraw-addresses",
		Short: "List saved withdrawal addresses",
		Long: `Retrieve the withdrawal addresses saved on the account, optionally
filtered by --asset and --method.

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-withdrawal-addresses`,
		Run: showWithdrawAddresses,
	}

	withdrawInfoCmd = &cobra.Command{
		Use:   "withdraw-info",
		Short: "Preview withdrawal fee and limit",
		Long: `Retrieve fee information about a potential withdrawal without executing
it. Required: --asset, --key, --amount.

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-withdrawal-information`,
		Run: showWithdrawInfo,
	}

	withdrawCmd = &cobra.Command{
		Use:   "withdraw",
		Short: "Withdraw funds to a saved withdrawal key",
		Long: `Make a withdrawal request to a pre-configured withdrawal key (address).
Required: --asset, --key, --amount.

WARNING: this moves funds off the exchange and is irreversible.

Docs Link: https://docs.kraken.com/api/docs/rest-api/withdraw-funds`,
		Run: doWithdraw,
	}

	withdrawStatusCmd = &cobra.Command{
		Use:   "withdraw-status",
		Short: "Show status of recent withdrawals",
		Long: `Retrieve information about recent withdrawals, optionally filtered by
--asset and --method.

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-status-recent-withdrawals`,
		Run: showWithdrawStatus,
	}

	withdrawCancelCmd = &cobra.Command{
		Use:   "withdraw-cancel",
		Short: "Cancel a recently requested withdrawal",
		Long: `Cancel a recently requested withdrawal, if it has not already been
processed. Required: --asset, --refid.

Docs Link: https://docs.kraken.com/api/docs/rest-api/cancel-withdrawal`,
		Run: doWithdrawCancel,
	}
)

// InitWithdrawCmds returns the withdrawal-related subcommands.
func InitWithdrawCmds() []*cobra.Command {
	withdrawMethodsCmd.Flags().StringP("asset", "a", "", "filter by asset")
	withdrawMethodsCmd.Flags().StringP("network", "n", "", "filter by network")

	withdrawAddressesCmd.Flags().StringP("asset", "a", "", "filter by asset")
	withdrawAddressesCmd.Flags().StringP("method", "m", "", "filter by method")

	withdrawInfoCmd.Flags().StringP("asset", "a", "", "asset to withdraw (required)")
	withdrawInfoCmd.Flags().StringP("key", "k", "", "withdrawal key name (required)")
	withdrawInfoCmd.Flags().StringP("amount", "m", "", "amount to withdraw (required)")
	withdrawInfoCmd.MarkFlagRequired("asset")
	withdrawInfoCmd.MarkFlagRequired("key")
	withdrawInfoCmd.MarkFlagRequired("amount")

	withdrawCmd.Flags().StringP("asset", "a", "", "asset to withdraw (required)")
	withdrawCmd.Flags().StringP("key", "k", "", "withdrawal key name (required)")
	withdrawCmd.Flags().StringP("amount", "m", "", "amount to withdraw (required)")
	withdrawCmd.Flags().String("address", "", "confirm the destination address matches the key")
	withdrawCmd.MarkFlagRequired("asset")
	withdrawCmd.MarkFlagRequired("key")
	withdrawCmd.MarkFlagRequired("amount")

	withdrawStatusCmd.Flags().StringP("asset", "a", "", "filter by asset")
	withdrawStatusCmd.Flags().StringP("method", "m", "", "filter by method")

	withdrawCancelCmd.Flags().StringP("asset", "a", "", "asset of the withdrawal (required)")
	withdrawCancelCmd.Flags().StringP("refid", "r", "", "reference id of the withdrawal (required)")
	withdrawCancelCmd.MarkFlagRequired("asset")
	withdrawCancelCmd.MarkFlagRequired("refid")

	return []*cobra.Command{
		withdrawMethodsCmd, withdrawAddressesCmd, withdrawInfoCmd,
		withdrawCmd, withdrawStatusCmd, withdrawCancelCmd,
	}
}

func showWithdrawMethods(cmd *cobra.Command, args []string) {
	client := newClient()
	asset, _ := cmd.Flags().GetString("asset")
	network, _ := cmd.Flags().GetString("network")
	methods, err := client.GetWithdrawMethods(asset, network)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(methods)
}

func showWithdrawAddresses(cmd *cobra.Command, args []string) {
	client := newClient()
	asset, _ := cmd.Flags().GetString("asset")
	method, _ := cmd.Flags().GetString("method")
	addrs, err := client.GetWithdrawAddresses(asset, method)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(addrs)
}

func showWithdrawInfo(cmd *cobra.Command, args []string) {
	client := newClient()
	asset, _ := cmd.Flags().GetString("asset")
	key, _ := cmd.Flags().GetString("key")
	amountRaw, _ := cmd.Flags().GetString("amount")
	amount, err := decimal.NewFromString(amountRaw)
	if err != nil {
		log.Fatalf("invalid amount: %v", err)
	}
	info, err := client.GetWithdrawInfo(asset, key, amount)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(info)
}

func doWithdraw(cmd *cobra.Command, args []string) {
	client := newClient()
	asset, _ := cmd.Flags().GetString("asset")
	key, _ := cmd.Flags().GetString("key")
	amountRaw, _ := cmd.Flags().GetString("amount")
	address, _ := cmd.Flags().GetString("address")
	amount, err := decimal.NewFromString(amountRaw)
	if err != nil {
		log.Fatalf("invalid amount: %v", err)
	}
	ref, err := client.Withdraw(asset, key, amount, address)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(ref)
}

func showWithdrawStatus(cmd *cobra.Command, args []string) {
	client := newClient()
	asset, _ := cmd.Flags().GetString("asset")
	method, _ := cmd.Flags().GetString("method")
	statuses, err := client.GetWithdrawStatus(asset, method)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(statuses)
}

func doWithdrawCancel(cmd *cobra.Command, args []string) {
	client := newClient()
	asset, _ := cmd.Flags().GetString("asset")
	refid, _ := cmd.Flags().GetString("refid")
	ok, err := client.WithdrawCancel(asset, refid)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(map[string]bool{"cancelled": ok})
}
