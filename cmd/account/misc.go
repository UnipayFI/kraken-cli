package account

import (
	"log"
	"strings"

	"github.com/UnipayFI/kraken-cli/printer"
	"github.com/spf13/cobra"
)

var (
	openPositionsCmd = &cobra.Command{
		Use:   "open-positions",
		Short: "List open margin positions",
		Long: `Retrieve information about currently open margin positions.

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-open-positions`,
		Run: showOpenPositions,
	}

	ledgersCmd = &cobra.Command{
		Use:   "ledgers",
		Short: "List ledger entries",
		Long: `Retrieve information about ledger entries. 50 results are returned at a
time, use offset for pagination.

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-ledgers-info`,
		Run: showLedgers,
	}

	queryLedgersCmd = &cobra.Command{
		Use:   "query-ledgers <id> [id...]",
		Short: "Query specific ledger entries by id (max 20)",
		Long: `Retrieve information about specific ledger entries by id.

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-ledgers`,
		Args: cobra.MinimumNArgs(1),
		Run:  showQueryLedgers,
	}

	tradeVolumeCmd = &cobra.Command{
		Use:   "trade-volume",
		Short: "Show 30-day trade volume and fee tier",
		Long: `Retrieve the account's 30-day USD trade volume and applicable fee tier.

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-trade-volume`,
		Run: showTradeVolume,
	}

	apiKeyInfoCmd = &cobra.Command{
		Use:   "api-key-info",
		Short: "Show information about the signing API key",
		Long: `Retrieve metadata (name, permissions, timestamps) about the API key used
to sign requests.

Docs Link: https://docs.kraken.com/api/docs/rest-api/get-api-key-info`,
		Run: showApiKeyInfo,
	}
)

// InitMiscCmds returns the position/ledger/volume/api-key subcommands.
func InitMiscCmds() []*cobra.Command {
	openPositionsCmd.Flags().BoolP("docalcs", "c", false, "include profit/loss calculations")

	ledgersCmd.Flags().StringP("asset", "s", "", "comma-separated assets to filter by")
	ledgersCmd.Flags().StringP("type", "T", "", "ledger type: all, deposit, withdrawal, trade, margin, ...")
	ledgersCmd.Flags().StringP("start", "a", "", "starting unix timestamp or ledger id")
	ledgersCmd.Flags().StringP("end", "e", "", "ending unix timestamp or ledger id")
	ledgersCmd.Flags().IntP("offset", "o", 0, "result offset for pagination")

	tradeVolumeCmd.Flags().StringP("pair", "p", "", "comma-separated pairs to include fee schedule for")

	return []*cobra.Command{openPositionsCmd, ledgersCmd, queryLedgersCmd, tradeVolumeCmd, apiKeyInfoCmd}
}

func splitCSV(v string) []string {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func showOpenPositions(cmd *cobra.Command, args []string) {
	client := newClient()
	docalcs, _ := cmd.Flags().GetBool("docalcs")
	positions, err := client.GetOpenPositions(docalcs)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(positions)
}

func showLedgers(cmd *cobra.Command, args []string) {
	client := newClient()
	asset, _ := cmd.Flags().GetString("asset")
	ledgerType, _ := cmd.Flags().GetString("type")
	start, _ := cmd.Flags().GetString("start")
	end, _ := cmd.Flags().GetString("end")
	offset, _ := cmd.Flags().GetInt("offset")
	ledgers, err := client.GetLedgers(splitCSV(asset), ledgerType, start, end, offset)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(ledgers)
}

func showQueryLedgers(cmd *cobra.Command, args []string) {
	client := newClient()
	ledgers, err := client.QueryLedgers(args)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(ledgers)
}

func showTradeVolume(cmd *cobra.Command, args []string) {
	client := newClient()
	pair, _ := cmd.Flags().GetString("pair")
	tv, err := client.GetTradeVolume(splitCSV(pair))
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(tv)
}

func showApiKeyInfo(cmd *cobra.Command, args []string) {
	client := newClient()
	info, err := client.GetApiKeyInfo()
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(info)
}
