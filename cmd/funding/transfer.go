package funding

import (
	"log"

	"github.com/UnipayFI/kraken-cli/printer"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var walletTransferCmd = &cobra.Command{
	Use:   "wallet-transfer",
	Short: "Transfer funds between Spot and Futures wallets",
	Long: `Transfer funds between the Spot and Futures wallets.
Required: --asset, --amount. Defaults transfer from "Spot Wallet" to
"Futures Wallet"; override with --from and --to.

Docs Link: https://docs.kraken.com/api/docs/rest-api/wallet-transfer`,
	Run: doWalletTransfer,
}

// InitTransferCmds returns the wallet-transfer subcommand.
func InitTransferCmds() []*cobra.Command {
	walletTransferCmd.Flags().StringP("asset", "a", "", "asset to transfer (required)")
	walletTransferCmd.Flags().StringP("amount", "m", "", "amount to transfer (required)")
	walletTransferCmd.Flags().StringP("from", "f", "", "source wallet (default \"Spot Wallet\")")
	walletTransferCmd.Flags().StringP("to", "t", "", "destination wallet (default \"Futures Wallet\")")
	walletTransferCmd.MarkFlagRequired("asset")
	walletTransferCmd.MarkFlagRequired("amount")
	return []*cobra.Command{walletTransferCmd}
}

func doWalletTransfer(cmd *cobra.Command, args []string) {
	client := newClient()
	asset, _ := cmd.Flags().GetString("asset")
	amountRaw, _ := cmd.Flags().GetString("amount")
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	amount, err := decimal.NewFromString(amountRaw)
	if err != nil {
		log.Fatalf("invalid amount: %v", err)
	}
	if !amount.IsPositive() {
		log.Fatal("amount must be greater than 0")
	}
	ref, err := client.WalletTransfer(asset, amount, from, to)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(ref)
}
