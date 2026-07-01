package order

import (
	"strings"

	kraken "github.com/UnipayFI/go-kraken"
	"github.com/UnipayFI/kraken-cli/common"
	"github.com/UnipayFI/kraken-cli/printer"
)

// ---------------------------------------------------------------------------
// Add order result
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*AddOrderResult)(nil)

type AddOrderResult struct {
	kraken.AddOrderResult
}

func (a *AddOrderResult) Header() []string {
	return []string{"TxID", "Order Description", "Close Description"}
}

func (a *AddOrderResult) Row() [][]any {
	txid := strings.Join(a.TxID, ", ")
	return [][]any{
		{txid, a.Description.Order, a.Description.Close},
	}
}

// ---------------------------------------------------------------------------
// Amend order result
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*AmendOrderResult)(nil)

type AmendOrderResult struct {
	kraken.AmendOrderResult
}

func (a *AmendOrderResult) Header() []string {
	return []string{"Amend ID"}
}

func (a *AmendOrderResult) Row() [][]any {
	return [][]any{{a.AmendID}}
}

// ---------------------------------------------------------------------------
// Edit order result
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*EditOrderResult)(nil)

type EditOrderResult struct {
	kraken.EditOrderResult
}

func (e *EditOrderResult) Header() []string {
	return []string{"Status", "TxID", "Original TxID", "Cancelled", "Volume", "Price", "Description"}
}

func (e *EditOrderResult) Row() [][]any {
	return [][]any{
		{e.Status, e.TxID, e.OriginalTxID, e.OrdersCancelled, e.Volume, e.Price, e.Description.Order},
	}
}

// ---------------------------------------------------------------------------
// Cancel results (count + pending)
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*CancelResult)(nil)

// CancelResult is the shared shape of the cancel endpoints: a count of orders
// affected plus, for single cancels, whether the cancellation is still pending.
type CancelResult struct {
	Count   int  `json:"count"`
	Pending bool `json:"pending"`
}

func (c *CancelResult) Header() []string {
	return []string{"Cancelled Count", "Pending"}
}

func (c *CancelResult) Row() [][]any {
	return [][]any{{c.Count, c.Pending}}
}

// ---------------------------------------------------------------------------
// Cancel-all-after result (dead man's switch)
// ---------------------------------------------------------------------------

var _ printer.TableWriter = (*CancelAllAfterResult)(nil)

type CancelAllAfterResult struct {
	kraken.CancelAllOrdersAfterResult
}

func (c *CancelAllAfterResult) Header() []string {
	return []string{"Current Time", "Trigger Time"}
}

func (c *CancelAllAfterResult) Row() [][]any {
	return [][]any{
		{common.FormatTime(c.CurrentTime), common.FormatTime(c.TriggerTime)},
	}
}
