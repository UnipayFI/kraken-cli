package account

import (
	"log"

	"github.com/UnipayFI/kraken-cli/printer"
	"github.com/spf13/cobra"
)

var (
	tradesHistoryCmd = &cobra.Command{
		Use:   "trades-history",
		Short: "List trade history",
		Long: `Retrieve information about trades/fills. 50 results are returned at a
time, use offset for pagination.

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-trade-history`,
		Run: showTradesHistory,
	}

	queryTradesCmd = &cobra.Command{
		Use:   "query-trades <txid> [txid...]",
		Short: "Query specific trades by tx id (max 20)",
		Long: `Retrieve information about specific trades/fills by transaction id.

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-trades-info`,
		Args: cobra.MinimumNArgs(1),
		Run:  showQueryTrades,
	}
)

// InitTradeCmds returns the trade-query subcommands.
func InitTradeCmds() []*cobra.Command {
	tradesHistoryCmd.Flags().StringP("type", "T", "", "trade type: all, any position, closed position, closing position, no position")
	tradesHistoryCmd.Flags().StringP("start", "a", "", "starting unix timestamp or trade tx id")
	tradesHistoryCmd.Flags().StringP("end", "e", "", "ending unix timestamp or trade tx id")
	tradesHistoryCmd.Flags().IntP("offset", "o", 0, "result offset for pagination")

	return []*cobra.Command{tradesHistoryCmd, queryTradesCmd}
}

func showTradesHistory(cmd *cobra.Command, args []string) {
	client := newClient()
	tradeType, _ := cmd.Flags().GetString("type")
	start, _ := cmd.Flags().GetString("start")
	end, _ := cmd.Flags().GetString("end")
	offset, _ := cmd.Flags().GetInt("offset")
	trades, err := client.GetTradesHistory(tradeType, start, end, offset)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(trades)
}

func showQueryTrades(cmd *cobra.Command, args []string) {
	client := newClient()
	trades, err := client.QueryTrades(args)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(trades)
}
