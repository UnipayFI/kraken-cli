package account

import (
	"log"

	"github.com/UnipayFI/kraken-cli/config"
	"github.com/UnipayFI/kraken-cli/exchange"
	"github.com/UnipayFI/kraken-cli/exchange/account"
	"github.com/UnipayFI/kraken-cli/printer"
	"github.com/spf13/cobra"
)

func newClient() account.Client {
	return account.Client{Client: exchange.NewClient(config.Config.APIKey, config.Config.APISecret)}
}

var (
	balanceCmd = &cobra.Command{
		Use:   "balance",
		Short: "Show account cash balances",
		Long: `Retrieve all cash balances, net of pending withdrawals (non-zero only).

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-account-balance`,
		Run: showBalance,
	}

	balanceExCmd = &cobra.Command{
		Use:   "balance-ex",
		Short: "Show extended balances (with amount held for orders)",
		Long: `Retrieve all cash balances with a breakdown of the amount held for open
orders/positions (non-zero only).

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-extended-balance`,
		Run: showBalanceEx,
	}

	tradeBalanceCmd = &cobra.Command{
		Use:   "trade-balance",
		Short: "Show margin/equity summary",
		Long: `Retrieve a summary of collateral balances, margin position valuations,
equity and margin level.

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-trade-balance`,
		Run: showTradeBalance,
	}
)

// InitBalanceCmds returns the balance-related subcommands.
func InitBalanceCmds() []*cobra.Command {
	tradeBalanceCmd.Flags().StringP("asset", "a", "", "base asset for the balance (default ZUSD)")
	return []*cobra.Command{balanceCmd, balanceExCmd, tradeBalanceCmd}
}

func showBalance(cmd *cobra.Command, args []string) {
	client := newClient()
	balances, err := client.GetBalance()
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(balances)
}

func showBalanceEx(cmd *cobra.Command, args []string) {
	client := newClient()
	balances, err := client.GetExtendedBalance()
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(balances)
}

func showTradeBalance(cmd *cobra.Command, args []string) {
	client := newClient()
	asset, _ := cmd.Flags().GetString("asset")
	tb, err := client.GetTradeBalance(asset)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(tb)
}
