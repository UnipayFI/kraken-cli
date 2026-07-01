package order

import (
	"log"

	"github.com/UnipayFI/kraken-cli/config"
	"github.com/UnipayFI/kraken-cli/exchange"
	"github.com/UnipayFI/kraken-cli/exchange/order"
	"github.com/UnipayFI/kraken-cli/printer"
	"github.com/spf13/cobra"
)

func newClient() order.Client {
	return order.Client{Client: exchange.NewClient(config.Config.APIKey, config.Config.APISecret)}
}

var (
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Place a new order",
		Long: `Place a new order.
Required: --pair, --side (buy/sell), --type, --volume.
Order types: market, limit, iceberg, stop-loss, take-profit, stop-loss-limit,
take-profit-limit, trailing-stop, trailing-stop-limit, settle-position.
Use --validate to validate inputs without submitting the order.

Docs Link: https://docs.kraken.com/api/docs/rest-api/add-order`,
		Run: addOrder,
	}

	amendCmd = &cobra.Command{
		Use:   "amend",
		Short: "Amend an open order in place",
		Long: `Amend an open order while keeping its queue priority where possible.
Identify the order with --txid or --client-order-id.

Docs Link: https://docs.kraken.com/api/docs/rest-api/amend-order`,
		Run: amendOrder,
	}

	editCmd = &cobra.Command{
		Use:   "edit",
		Short: "Cancel-replace an order (new tx id)",
		Long: `Edit the volume and/or price of an open order. This cancels the original
order and places a new one with a new tx id.
Required: --pair, --txid.

Docs Link: https://docs.kraken.com/api/docs/rest-api/edit-order`,
		Run: editOrder,
	}

	cancelCmd = &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a single open order",
		Long: `Cancel a particular open order (or set of open orders) by --txid or
--client-order-id.

Docs Link: https://docs.kraken.com/api/docs/rest-api/cancel-order`,
		Run: cancelOrder,
	}

	cancelAllCmd = &cobra.Command{
		Use:   "cancel-all",
		Short: "Cancel all open orders",
		Long: `Cancel all open orders.

Docs Link: https://docs.kraken.com/api/docs/rest-api/cancel-all-orders`,
		Run: cancelAll,
	}

	cancelAllAfterCmd = &cobra.Command{
		Use:   "cancel-all-after",
		Short: "Arm the dead man's switch",
		Long: `Set (or, with --timeout 0, disarm) a countdown that cancels all orders
after the given number of seconds unless the timer is repeatedly extended. Meant
as protection against network failures.

Docs Link: https://docs.kraken.com/api/docs/rest-api/cancel-all-orders-after`,
		Run: cancelAllAfter,
	}

	cancelBatchCmd = &cobra.Command{
		Use:   "cancel-batch <txid> [txid...]",
		Short: "Cancel multiple orders by tx id (max 50)",
		Long: `Cancel multiple open orders by transaction id or user reference.

Docs Link: https://docs.kraken.com/api/docs/rest-api/cancel-order-batch`,
		Args: cobra.MinimumNArgs(1),
		Run:  cancelBatch,
	}
)

// InitCmds returns all order (trading) subcommands.
func InitCmds() []*cobra.Command {
	addCmd.Flags().StringP("pair", "p", "", "asset pair, e.g. XBTUSD (required)")
	addCmd.Flags().StringP("side", "S", "", "buy or sell (required)")
	addCmd.Flags().StringP("type", "t", "", "order type, e.g. limit or market (required)")
	addCmd.Flags().StringP("volume", "v", "", "order volume in base currency (required)")
	addCmd.Flags().StringP("price", "P", "", "limit price / trigger price")
	addCmd.Flags().String("price2", "", "secondary price (for stop-loss-limit etc.)")
	addCmd.Flags().String("leverage", "", "leverage, e.g. 5:1")
	addCmd.Flags().String("oflags", "", "comma-separated order flags: post, fcib, fciq, nompp, viqc")
	addCmd.Flags().StringP("tif", "T", "", "time in force: GTC, IOC, GTD, FOK")
	addCmd.Flags().String("trigger", "", "trigger signal: last or index")
	addCmd.Flags().String("stptype", "", "self-trade prevention: cancel-newest, cancel-oldest, cancel-both")
	addCmd.Flags().Bool("reduce-only", false, "reduce-only order")
	addCmd.Flags().String("display-volume", "", "display volume for iceberg orders")
	addCmd.Flags().String("client-order-id", "", "client order id")
	addCmd.Flags().String("start-time", "", "scheduled start time (0, unix ts, or +seconds)")
	addCmd.Flags().String("expire-time", "", "expire time for GTD (unix ts or +seconds)")
	addCmd.Flags().Bool("validate", false, "validate inputs only, do not submit")
	addCmd.MarkFlagRequired("pair")
	addCmd.MarkFlagRequired("side")
	addCmd.MarkFlagRequired("type")
	addCmd.MarkFlagRequired("volume")

	amendCmd.Flags().StringP("txid", "i", "", "order tx id to amend")
	amendCmd.Flags().StringP("client-order-id", "c", "", "client order id to amend")
	amendCmd.Flags().StringP("volume", "v", "", "new order quantity")
	amendCmd.Flags().String("display-volume", "", "new display quantity (iceberg)")
	amendCmd.Flags().StringP("price", "P", "", "new limit price")
	amendCmd.Flags().String("trigger-price", "", "new trigger price")
	amendCmd.Flags().Bool("post-only", false, "restrict to post-only")

	editCmd.Flags().StringP("pair", "p", "", "asset pair (required)")
	editCmd.Flags().StringP("txid", "i", "", "order tx id to edit (required)")
	editCmd.Flags().StringP("volume", "v", "", "new order volume")
	editCmd.Flags().StringP("price", "P", "", "new price")
	editCmd.Flags().String("price2", "", "new secondary price")
	editCmd.Flags().String("oflags", "", "comma-separated order flags")
	editCmd.MarkFlagRequired("pair")
	editCmd.MarkFlagRequired("txid")

	cancelCmd.Flags().StringP("txid", "i", "", "order tx id or user reference to cancel")
	cancelCmd.Flags().StringP("client-order-id", "c", "", "client order id to cancel")

	cancelAllAfterCmd.Flags().IntP("timeout", "t", 0, "timeout in seconds (0 disarms)")
	cancelAllAfterCmd.MarkFlagRequired("timeout")

	return []*cobra.Command{addCmd, amendCmd, editCmd, cancelCmd, cancelAllCmd, cancelAllAfterCmd, cancelBatchCmd}
}

func addOrder(cmd *cobra.Command, args []string) {
	client := newClient()
	params := map[string]string{}
	strFlags := map[string]string{
		"pair":            "pair",
		"side":            "side",
		"type":            "type",
		"volume":          "volume",
		"price":           "price",
		"price2":          "price2",
		"leverage":        "leverage",
		"oflags":          "oflags",
		"tif":             "timeInForce",
		"trigger":         "trigger",
		"stptype":         "stptype",
		"display-volume":  "displayVolume",
		"client-order-id": "clientOrderId",
		"start-time":      "startTime",
		"expire-time":     "expireTime",
	}
	for flag, key := range strFlags {
		if v, _ := cmd.Flags().GetString(flag); v != "" {
			params[key] = v
		}
	}
	if v, _ := cmd.Flags().GetBool("reduce-only"); v {
		params["reduceOnly"] = "true"
	}
	if v, _ := cmd.Flags().GetBool("validate"); v {
		params["validate"] = "true"
	}

	res, err := client.AddOrder(params)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(res)
}

func amendOrder(cmd *cobra.Command, args []string) {
	client := newClient()
	txid, _ := cmd.Flags().GetString("txid")
	clientOrderID, _ := cmd.Flags().GetString("client-order-id")
	if txid == "" && clientOrderID == "" {
		log.Fatal("either --txid or --client-order-id is required")
	}
	params := map[string]string{}
	for flag, key := range map[string]string{
		"volume":         "volume",
		"display-volume": "displayVolume",
		"price":          "price",
		"trigger-price":  "triggerPrice",
	} {
		if v, _ := cmd.Flags().GetString(flag); v != "" {
			params[key] = v
		}
	}
	if v, _ := cmd.Flags().GetBool("post-only"); v {
		params["postOnly"] = "true"
	}

	res, err := client.AmendOrder(txid, clientOrderID, params)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(res)
}

func editOrder(cmd *cobra.Command, args []string) {
	client := newClient()
	pair, _ := cmd.Flags().GetString("pair")
	txid, _ := cmd.Flags().GetString("txid")
	params := map[string]string{}
	for flag, key := range map[string]string{
		"volume": "volume",
		"price":  "price",
		"price2": "price2",
		"oflags": "oflags",
	} {
		if v, _ := cmd.Flags().GetString(flag); v != "" {
			params[key] = v
		}
	}

	res, err := client.EditOrder(pair, txid, params)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(res)
}

func cancelOrder(cmd *cobra.Command, args []string) {
	client := newClient()
	txid, _ := cmd.Flags().GetString("txid")
	clientOrderID, _ := cmd.Flags().GetString("client-order-id")
	if txid == "" && clientOrderID == "" {
		log.Fatal("either --txid or --client-order-id is required")
	}
	res, err := client.CancelOrder(txid, clientOrderID)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(res)
}

func cancelAll(cmd *cobra.Command, args []string) {
	client := newClient()
	res, err := client.CancelAll()
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(res)
}

func cancelAllAfter(cmd *cobra.Command, args []string) {
	client := newClient()
	timeout, _ := cmd.Flags().GetInt("timeout")
	res, err := client.CancelAllAfter(timeout)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(res)
}

func cancelBatch(cmd *cobra.Command, args []string) {
	client := newClient()
	res, err := client.CancelBatch(args)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(res)
}
