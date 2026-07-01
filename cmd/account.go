package cmd

import (
	"github.com/UnipayFI/kraken-cli/cmd/account"
	"github.com/spf13/cobra"
)

var AccountCmd = &cobra.Command{
	Use:   "account",
	Short: "Account data commands (balances, orders, trades, ledgers)",
}

func init() {
	AccountCmd.AddCommand(account.InitBalanceCmds()...)
	AccountCmd.AddCommand(account.InitOrderCmds()...)
	AccountCmd.AddCommand(account.InitTradeCmds()...)
	AccountCmd.AddCommand(account.InitMiscCmds()...)
	RootCmd.AddCommand(AccountCmd)
}
