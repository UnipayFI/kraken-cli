package cmd

import (
	"github.com/UnipayFI/kraken-cli/cmd/funding"
	"github.com/spf13/cobra"
)

var FundingCmd = &cobra.Command{
	Use:   "funding",
	Short: "Funding commands (deposits, withdrawals, transfers)",
}

func init() {
	FundingCmd.AddCommand(funding.InitDepositCmds()...)
	FundingCmd.AddCommand(funding.InitWithdrawCmds()...)
	FundingCmd.AddCommand(funding.InitTransferCmds()...)
	RootCmd.AddCommand(FundingCmd)
}
