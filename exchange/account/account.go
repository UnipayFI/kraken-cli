package account

import (
	"context"

	"github.com/UnipayFI/kraken-cli/exchange"
)

// Client wraps the shared exchange client with account-scoped query methods.
type Client struct {
	*exchange.Client
}

func ctx() context.Context { return context.Background() }

// GetBalance returns non-zero cash balances sorted by asset.
func (c *Client) GetBalance() (*BalanceList, error) {
	m, err := c.NewKrakenClient().NewGetAccountBalanceService().Do(ctx())
	if err != nil {
		return nil, err
	}
	return newBalanceList(m), nil
}

// GetExtendedBalance returns non-zero balances with the amount held for orders.
func (c *Client) GetExtendedBalance() (*ExtendedBalanceList, error) {
	m, err := c.NewKrakenClient().NewGetExtendedBalanceService().Do(ctx())
	if err != nil {
		return nil, err
	}
	return newExtendedBalanceList(m), nil
}

// GetTradeBalance returns the margin/equity summary for the given base asset
// (empty asset uses the Kraken default).
func (c *Client) GetTradeBalance(asset string) (*TradeBalance, error) {
	service := c.NewKrakenClient().NewGetTradeBalanceService()
	if asset != "" {
		service.SetAsset(asset)
	}
	tb, err := service.Do(ctx())
	if err != nil {
		return nil, err
	}
	return &TradeBalance{TradeBalance: *tb}, nil
}

// GetOpenOrders returns the account's open orders.
func (c *Client) GetOpenOrders(trades bool) (*OrderList, error) {
	service := c.NewKrakenClient().NewGetOpenOrdersService()
	if trades {
		service.SetTrades(true)
	}
	res, err := service.Do(ctx())
	if err != nil {
		return nil, err
	}
	return newOrderList(res.Open), nil
}

// GetClosedOrders returns the account's closed orders with optional time range.
func (c *Client) GetClosedOrders(start, end string, offset int, trades bool) (*OrderList, error) {
	service := c.NewKrakenClient().NewGetClosedOrdersService()
	if start != "" {
		service.SetStart(start)
	}
	if end != "" {
		service.SetEnd(end)
	}
	if offset != 0 {
		service.SetOffset(offset)
	}
	if trades {
		service.SetTrades(true)
	}
	res, err := service.Do(ctx())
	if err != nil {
		return nil, err
	}
	return newOrderList(res.Closed), nil
}

// QueryOrders returns full info for specific order tx ids (max 50).
func (c *Client) QueryOrders(txids []string, trades bool) (*OrderList, error) {
	service := c.NewKrakenClient().NewQueryOrdersService(txids...)
	if trades {
		service.SetTrades(true)
	}
	m, err := service.Do(ctx())
	if err != nil {
		return nil, err
	}
	return newOrderList(m), nil
}

// GetTradesHistory returns the account's trade history with optional range.
func (c *Client) GetTradesHistory(tradeType, start, end string, offset int) (*TradeList, error) {
	service := c.NewKrakenClient().NewGetTradesHistoryService()
	if tradeType != "" {
		service.SetType(tradeType)
	}
	if start != "" {
		service.SetStart(start)
	}
	if end != "" {
		service.SetEnd(end)
	}
	if offset != 0 {
		service.SetOffset(offset)
	}
	res, err := service.Do(ctx())
	if err != nil {
		return nil, err
	}
	return newTradeList(res.Trades), nil
}

// QueryTrades returns full info for specific trade tx ids (max 20).
func (c *Client) QueryTrades(txids []string) (*TradeList, error) {
	m, err := c.NewKrakenClient().NewQueryTradesService(txids...).Do(ctx())
	if err != nil {
		return nil, err
	}
	return newTradeList(m), nil
}

// GetOpenPositions returns the account's open margin positions.
func (c *Client) GetOpenPositions(doCalcs bool) (*PositionList, error) {
	service := c.NewKrakenClient().NewGetOpenPositionsService()
	if doCalcs {
		service.SetDoCalcs(true)
	}
	m, err := service.Do(ctx())
	if err != nil {
		return nil, err
	}
	return newPositionList(m), nil
}

// GetLedgers returns the account's ledger entries with optional filters.
func (c *Client) GetLedgers(assets []string, ledgerType, start, end string, offset int) (*LedgerList, error) {
	service := c.NewKrakenClient().NewGetLedgersService()
	if len(assets) > 0 {
		service.SetAsset(assets...)
	}
	if ledgerType != "" {
		service.SetType(ledgerType)
	}
	if start != "" {
		service.SetStart(start)
	}
	if end != "" {
		service.SetEnd(end)
	}
	if offset != 0 {
		service.SetOffset(offset)
	}
	res, err := service.Do(ctx())
	if err != nil {
		return nil, err
	}
	return newLedgerList(res.Ledger), nil
}

// QueryLedgers returns full info for specific ledger ids (max 20).
func (c *Client) QueryLedgers(ids []string) (*LedgerList, error) {
	m, err := c.NewKrakenClient().NewQueryLedgersService(ids...).Do(ctx())
	if err != nil {
		return nil, err
	}
	return newLedgerList(m), nil
}

// GetTradeVolume returns the account's 30-day trade volume and fee tier.
func (c *Client) GetTradeVolume(pairs []string) (*TradeVolume, error) {
	service := c.NewKrakenClient().NewGetTradeVolumeService()
	if len(pairs) > 0 {
		service.SetPair(pairs...)
	}
	tv, err := service.Do(ctx())
	if err != nil {
		return nil, err
	}
	return &TradeVolume{TradeVolume: *tv}, nil
}

// GetApiKeyInfo returns metadata about the signing API key.
func (c *Client) GetApiKeyInfo() (*APIKeyInfo, error) {
	info, err := c.NewKrakenClient().NewGetApiKeyInfoService().Do(ctx())
	if err != nil {
		return nil, err
	}
	return &APIKeyInfo{APIKeyInfo: *info}, nil
}
