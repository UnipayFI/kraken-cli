package account

import (
	"sort"
	"strings"

	kraken "github.com/UnipayFI/go-kraken"
	"github.com/UnipayFI/kraken-cli/common"
	"github.com/UnipayFI/kraken-cli/printer"
	"github.com/shopspring/decimal"
)

// ---------------------------------------------------------------------------
// Balances
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*BalanceList)(nil)

// Balance is one asset's cash balance.
type Balance struct {
	Asset   string          `json:"asset"`
	Balance decimal.Decimal `json:"balance"`
}

type BalanceList []Balance

func (b *BalanceList) Header() []string {
	return []string{"Asset", "Balance"}
}

func (b *BalanceList) Row() [][]any {
	rows := [][]any{}
	for _, e := range *b {
		rows = append(rows, []any{e.Asset, e.Balance})
	}
	return rows
}

// newBalanceList builds a list sorted by asset from the raw balance map,
// dropping zero balances.
func newBalanceList(m map[string]decimal.Decimal) *BalanceList {
	list := BalanceList{}
	for asset, amount := range m {
		if amount.IsZero() {
			continue
		}
		list = append(list, Balance{Asset: asset, Balance: amount})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Asset < list[j].Asset })
	return &list
}

var _ printer.TableWriter = (*ExtendedBalanceList)(nil)

// ExtendedBalance is one asset's total balance plus the portion held for orders.
type ExtendedBalance struct {
	Asset     string          `json:"asset"`
	Balance   decimal.Decimal `json:"balance"`
	HoldTrade decimal.Decimal `json:"hold_trade"`
}

type ExtendedBalanceList []ExtendedBalance

func (b *ExtendedBalanceList) Header() []string {
	return []string{"Asset", "Balance", "Hold Trade"}
}

func (b *ExtendedBalanceList) Row() [][]any {
	rows := [][]any{}
	for _, e := range *b {
		rows = append(rows, []any{e.Asset, e.Balance, e.HoldTrade})
	}
	return rows
}

func newExtendedBalanceList(m map[string]kraken.ExtendedBalance) *ExtendedBalanceList {
	list := ExtendedBalanceList{}
	for asset, b := range m {
		if b.Balance.IsZero() && b.HoldTrade.IsZero() {
			continue
		}
		list = append(list, ExtendedBalance{Asset: asset, Balance: b.Balance, HoldTrade: b.HoldTrade})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Asset < list[j].Asset })
	return &list
}

// ---------------------------------------------------------------------------
// Trade balance
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*TradeBalance)(nil)

type TradeBalance struct {
	kraken.TradeBalance
}

func (t *TradeBalance) Header() []string {
	return []string{"Equiv Balance", "Trade Balance", "Equity", "Free Margin", "Margin Level", "Unrealized PnL", "Margin (Open)"}
}

func (t *TradeBalance) Row() [][]any {
	return [][]any{
		{t.EquivalentBalance, t.TradeBalance.TradeBalance, t.Equity, t.FreeMargin, t.MarginLevel, t.UnrealizedNetPnL, t.MarginOpenPositions},
	}
}

// ---------------------------------------------------------------------------
// Orders
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*OrderList)(nil)

// Order carries an order's tx id alongside its info (the API keys the map by
// tx id).
type Order struct {
	TxID string `json:"txid"`
	kraken.OrderInfo
}

type OrderList []Order

func (o *OrderList) Header() []string {
	return []string{"TxID", "Pair", "Side", "Type", "Status", "Price", "Avg Price", "Volume", "Executed", "Cost", "Fee", "Open Time"}
}

func (o *OrderList) Row() [][]any {
	rows := [][]any{}
	for _, ord := range *o {
		rows = append(rows, []any{
			ord.TxID,
			ord.Description.Pair,
			ord.Description.Type,
			ord.Description.OrderType,
			ord.Status,
			ord.Description.Price,
			ord.Price,
			ord.Volume,
			ord.VolumeExecuted,
			ord.Cost,
			ord.Fee,
			common.FormatTime(ord.OpenTime),
		})
	}
	return rows
}

func newOrderList(m map[string]kraken.OrderInfo) *OrderList {
	list := OrderList{}
	for txid, info := range m {
		list = append(list, Order{TxID: txid, OrderInfo: info})
	}
	// Most-recent first.
	sort.Slice(list, func(i, j int) bool { return list[i].OpenTime.After(list[j].OpenTime) })
	return &list
}

// ---------------------------------------------------------------------------
// Trades
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*TradeList)(nil)

type Trade struct {
	TxID string `json:"txid"`
	kraken.TradeHistoryEntry
}

type TradeList []Trade

func (t *TradeList) Header() []string {
	return []string{"TxID", "Order TxID", "Pair", "Side", "Type", "Price", "Cost", "Fee", "Volume", "Maker", "Time"}
}

func (t *TradeList) Row() [][]any {
	rows := [][]any{}
	for _, tr := range *t {
		rows = append(rows, []any{
			tr.TxID,
			tr.OrderTxID,
			tr.Pair,
			tr.Type,
			tr.OrderType,
			tr.Price,
			tr.Cost,
			tr.Fee,
			tr.Volume,
			tr.Maker,
			common.FormatTime(tr.Time),
		})
	}
	return rows
}

func newTradeList(m map[string]kraken.TradeHistoryEntry) *TradeList {
	list := TradeList{}
	for txid, e := range m {
		list = append(list, Trade{TxID: txid, TradeHistoryEntry: e})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Time.After(list[j].Time) })
	return &list
}

// ---------------------------------------------------------------------------
// Positions
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*PositionList)(nil)

type Position struct {
	TxID string `json:"txid"`
	kraken.PositionInfo
}

type PositionList []Position

func (p *PositionList) Header() []string {
	return []string{"TxID", "Pair", "Side", "Type", "Volume", "Closed", "Cost", "Fee", "Margin", "Value", "Net", "Time"}
}

func (p *PositionList) Row() [][]any {
	rows := [][]any{}
	for _, pos := range *p {
		rows = append(rows, []any{
			pos.TxID,
			pos.Pair,
			pos.Type,
			pos.OrderType,
			pos.Volume,
			pos.VolumeClosed,
			pos.Cost,
			pos.Fee,
			pos.Margin,
			pos.Value,
			pos.Net,
			common.FormatTime(pos.Time),
		})
	}
	return rows
}

func newPositionList(m map[string]kraken.PositionInfo) *PositionList {
	list := PositionList{}
	for txid, info := range m {
		list = append(list, Position{TxID: txid, PositionInfo: info})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Time.After(list[j].Time) })
	return &list
}

// ---------------------------------------------------------------------------
// Ledgers
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*LedgerList)(nil)

type Ledger struct {
	ID string `json:"id"`
	kraken.LedgerEntry
}

type LedgerList []Ledger

func (l *LedgerList) Header() []string {
	return []string{"ID", "RefID", "Type", "SubType", "Asset", "Amount", "Fee", "Balance", "Time"}
}

func (l *LedgerList) Row() [][]any {
	rows := [][]any{}
	for _, e := range *l {
		rows = append(rows, []any{
			e.ID,
			e.RefID,
			e.Type,
			e.SubType,
			e.Asset,
			e.Amount,
			e.Fee,
			e.Balance,
			common.FormatTime(e.Time),
		})
	}
	return rows
}

func newLedgerList(m map[string]kraken.LedgerEntry) *LedgerList {
	list := LedgerList{}
	for id, e := range m {
		list = append(list, Ledger{ID: id, LedgerEntry: e})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Time.After(list[j].Time) })
	return &list
}

// ---------------------------------------------------------------------------
// Trade volume
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*TradeVolume)(nil)

type TradeVolume struct {
	kraken.TradeVolume
}

func (t *TradeVolume) Header() []string {
	return []string{"Currency", "30d Volume", "Asset Class"}
}

func (t *TradeVolume) Row() [][]any {
	return [][]any{
		{t.Currency, t.Volume, t.AssetClass},
	}
}

// ---------------------------------------------------------------------------
// API key info
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*APIKeyInfo)(nil)

type APIKeyInfo struct {
	kraken.APIKeyInfo
}

func (a *APIKeyInfo) Header() []string {
	return []string{"Key Name", "Permissions", "Created", "Last Used", "Valid Until"}
}

func (a *APIKeyInfo) Row() [][]any {
	validUntil := common.FormatTime(a.ValidUntil)
	if validUntil == "" {
		validUntil = "never"
	}
	perms := strings.Join(a.Permissions, ", ")
	return [][]any{
		{a.APIKeyName, perms, common.FormatTime(a.CreatedTime), common.FormatTime(a.LastUsed), validUntil},
	}
}
