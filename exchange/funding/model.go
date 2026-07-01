package funding

import (
	kraken "github.com/UnipayFI/go-kraken"
	"github.com/UnipayFI/kraken-cli/common"
	"github.com/UnipayFI/kraken-cli/printer"
)

// ---------------------------------------------------------------------------
// Deposit methods
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*DepositMethodList)(nil)

type DepositMethodList []kraken.DepositMethod

func (d *DepositMethodList) Header() []string {
	return []string{"Method", "Fee", "Minimum", "Gen Address", "Limit"}
}

func (d *DepositMethodList) Row() [][]any {
	rows := [][]any{}
	for _, m := range *d {
		rows = append(rows, []any{m.Method, m.Fee, m.Minimum, m.GenAddress, formatLimit(m.Limit)})
	}
	return rows
}

func formatLimit(l kraken.MethodLimit) string {
	if !l.HasLimit {
		return "none"
	}
	if l.Amount.IsZero() {
		return "yes"
	}
	return l.Amount.String()
}

// ---------------------------------------------------------------------------
// Deposit addresses
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*DepositAddressList)(nil)

type DepositAddressList []kraken.DepositAddress

func (d *DepositAddressList) Header() []string {
	return []string{"Address", "Tag", "Memo", "New", "Expire Time"}
}

func (d *DepositAddressList) Row() [][]any {
	rows := [][]any{}
	for _, a := range *d {
		rows = append(rows, []any{a.Address, a.Tag, a.Memo, a.New, common.FormatTime(a.ExpireTime)})
	}
	return rows
}

// ---------------------------------------------------------------------------
// Transfer status (deposits & withdrawals share the shape)
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*TransferStatusList)(nil)

type TransferStatusList []kraken.TransferStatus

func (t *TransferStatusList) Header() []string {
	return []string{"RefID", "Method", "Asset", "Amount", "Fee", "Status", "Time", "TxID"}
}

func (t *TransferStatusList) Row() [][]any {
	rows := [][]any{}
	for _, s := range *t {
		status := s.Status
		if s.StatusProp != "" {
			status = status + " (" + s.StatusProp + ")"
		}
		rows = append(rows, []any{s.RefID, s.Method, s.Asset, s.Amount, s.Fee, status, common.FormatTime(s.Time), s.TxID})
	}
	return rows
}

// ---------------------------------------------------------------------------
// Withdraw methods
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*WithdrawMethodList)(nil)

type WithdrawMethodList []kraken.WithdrawMethod

func (w *WithdrawMethodList) Header() []string {
	return []string{"Asset", "Method", "Network", "Minimum", "Fee", "Fee Asset"}
}

func (w *WithdrawMethodList) Row() [][]any {
	rows := [][]any{}
	for _, m := range *w {
		rows = append(rows, []any{m.Asset, m.Method, m.Network, m.Minimum, m.Fee.Fee, m.Fee.Asset})
	}
	return rows
}

// ---------------------------------------------------------------------------
// Withdraw addresses
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*WithdrawAddressList)(nil)

type WithdrawAddressList []kraken.WithdrawAddress

func (w *WithdrawAddressList) Header() []string {
	return []string{"Key", "Asset", "Method", "Network", "Address", "Verified"}
}

func (w *WithdrawAddressList) Row() [][]any {
	rows := [][]any{}
	for _, a := range *w {
		rows = append(rows, []any{a.Key, a.Asset, a.Method, a.Network, a.Address, a.Verified})
	}
	return rows
}

// ---------------------------------------------------------------------------
// Withdraw info (preview)
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*WithdrawInfo)(nil)

type WithdrawInfo struct {
	kraken.WithdrawInfo
}

func (w *WithdrawInfo) Header() []string {
	return []string{"Method", "Limit", "Amount (net)", "Fee"}
}

func (w *WithdrawInfo) Row() [][]any {
	return [][]any{{w.Method, w.Limit, w.Amount, w.Fee}}
}

// ---------------------------------------------------------------------------
// Reference-id results (Withdraw / WalletTransfer)
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*RefResult)(nil)

// RefResult wraps a single reference-id response.
type RefResult struct {
	RefID string `json:"refid"`
}

func (r *RefResult) Header() []string {
	return []string{"RefID"}
}

func (r *RefResult) Row() [][]any {
	return [][]any{{r.RefID}}
}
