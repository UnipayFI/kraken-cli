package account

import (
	"log"

	"github.com/UnipayFI/kraken-cli/printer"
	"github.com/spf13/cobra"
)

var (
	openOrdersCmd = &cobra.Command{
		Use:   "open-orders",
		Short: "List open orders",
		Long: `Retrieve information about currently open orders.

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-open-orders`,
		Run: showOpenOrders,
	}

	closedOrdersCmd = &cobra.Command{
		Use:   "closed-orders",
		Short: "List closed orders",
		Long: `Retrieve information about orders that have been closed (filled or
cancelled). 50 results are returned at a time, use offset for pagination.

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-closed-orders`,
		Run: showClosedOrders,
	}

	queryOrdersCmd = &cobra.Command{
		Use:   "query-orders <txid> [txid...]",
		Short: "Query specific orders by tx id (max 50)",
		Long: `Retrieve information about specific orders by transaction id.

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-orders-info`,
		Args: cobra.MinimumNArgs(1),
		Run:  showQueryOrders,
	}
)

// InitOrderCmds returns the order-query subcommands.
func InitOrderCmds() []*cobra.Command {
	openOrdersCmd.Flags().BoolP("trades", "t", false, "include related trade ids")

	closedOrdersCmd.Flags().BoolP("trades", "t", false, "include related trade ids")
	closedOrdersCmd.Flags().StringP("start", "a", "", "starting unix timestamp or order tx id")
	closedOrdersCmd.Flags().StringP("end", "e", "", "ending unix timestamp or order tx id")
	closedOrdersCmd.Flags().IntP("offset", "o", 0, "result offset for pagination")

	queryOrdersCmd.Flags().BoolP("trades", "t", false, "include related trade ids")

	return []*cobra.Command{openOrdersCmd, closedOrdersCmd, queryOrdersCmd}
}

func showOpenOrders(cmd *cobra.Command, args []string) {
	client := newClient()
	trades, _ := cmd.Flags().GetBool("trades")
	orders, err := client.GetOpenOrders(trades)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(orders)
}

func showClosedOrders(cmd *cobra.Command, args []string) {
	client := newClient()
	trades, _ := cmd.Flags().GetBool("trades")
	start, _ := cmd.Flags().GetString("start")
	end, _ := cmd.Flags().GetString("end")
	offset, _ := cmd.Flags().GetInt("offset")
	orders, err := client.GetClosedOrders(start, end, offset, trades)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(orders)
}

func showQueryOrders(cmd *cobra.Command, args []string) {
	client := newClient()
	trades, _ := cmd.Flags().GetBool("trades")
	orders, err := client.QueryOrders(args, trades)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(orders)
}
