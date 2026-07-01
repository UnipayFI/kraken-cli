package cmd

import (
	"github.com/UnipayFI/kraken-cli/cmd/order"
	"github.com/spf13/cobra"
)

var OrderCmd = &cobra.Command{
	Use:   "order",
	Short: "Trading commands (add, amend, edit, cancel orders)",
}

func init() {
	OrderCmd.AddCommand(order.InitCmds()...)
	RootCmd.AddCommand(OrderCmd)
}
