package order

import (
	"context"
	"strings"

	kraken "github.com/UnipayFI/go-kraken"
	"github.com/UnipayFI/kraken-cli/exchange"
	"github.com/shopspring/decimal"
)

// Client wraps the shared exchange client with trading methods.
type Client struct {
	*exchange.Client
}

func ctx() context.Context { return context.Background() }

// AddOrder places a new order. Required params: pair, side, type, volume.
// Optional: price, price2, leverage, oflags, timeInForce, trigger, stptype,
// reduceOnly, displayVolume, userref, clientOrderId, startTime, expireTime,
// validate.
func (c *Client) AddOrder(params map[string]string) (*AddOrderResult, error) {
	volume, err := decimal.NewFromString(params["volume"])
	if err != nil {
		return nil, err
	}
	side := kraken.OrderSide(strings.ToLower(params["side"]))
	orderType := kraken.OrderType(strings.ToLower(params["type"]))
	service := c.NewKrakenClient().NewAddOrderService(params["pair"], side, orderType, volume)

	if v := params["price"]; v != "" {
		p, err := decimal.NewFromString(v)
		if err != nil {
			return nil, err
		}
		service.SetPrice(p)
	}
	if v := params["price2"]; v != "" {
		p, err := decimal.NewFromString(v)
		if err != nil {
			return nil, err
		}
		service.SetPrice2(p)
	}
	if v := params["leverage"]; v != "" {
		service.SetLeverage(v)
	}
	if v := params["oflags"]; v != "" {
		service.SetOrderFlags(v)
	}
	if v := params["timeInForce"]; v != "" {
		service.SetTimeInForce(kraken.TimeInForce(strings.ToUpper(v)))
	}
	if v := params["trigger"]; v != "" {
		service.SetTrigger(kraken.TriggerSignal(strings.ToLower(v)))
	}
	if v := params["stptype"]; v != "" {
		service.SetSTPType(kraken.SelfTradePrevention(strings.ToLower(v)))
	}
	if params["reduceOnly"] == "true" {
		service.SetReduceOnly(true)
	}
	if v := params["displayVolume"]; v != "" {
		p, err := decimal.NewFromString(v)
		if err != nil {
			return nil, err
		}
		service.SetDisplayVolume(p)
	}
	if v := params["clientOrderId"]; v != "" {
		service.SetClientOrderID(v)
	}
	if v := params["startTime"]; v != "" {
		service.SetStartTime(v)
	}
	if v := params["expireTime"]; v != "" {
		service.SetExpireTime(v)
	}
	if params["validate"] == "true" {
		service.SetValidate(true)
	}

	res, err := service.Do(ctx())
	if err != nil {
		return nil, err
	}
	return &AddOrderResult{AddOrderResult: *res}, nil
}

// AmendOrder amends an open order in place (preserving queue priority). Either
// txid or clientOrderId identifies the order.
func (c *Client) AmendOrder(txid, clientOrderID string, params map[string]string) (*AmendOrderResult, error) {
	service := c.NewKrakenClient().NewAmendOrderService(txid)
	if clientOrderID != "" {
		service.SetClientOrderID(clientOrderID)
	}
	if v := params["volume"]; v != "" {
		p, err := decimal.NewFromString(v)
		if err != nil {
			return nil, err
		}
		service.SetOrderQty(p)
	}
	if v := params["displayVolume"]; v != "" {
		p, err := decimal.NewFromString(v)
		if err != nil {
			return nil, err
		}
		service.SetDisplayQty(p)
	}
	if v := params["price"]; v != "" {
		p, err := decimal.NewFromString(v)
		if err != nil {
			return nil, err
		}
		service.SetLimitPrice(p)
	}
	if v := params["triggerPrice"]; v != "" {
		p, err := decimal.NewFromString(v)
		if err != nil {
			return nil, err
		}
		service.SetTriggerPrice(p)
	}
	if params["postOnly"] == "true" {
		service.SetPostOnly(true)
	}
	res, err := service.Do(ctx())
	if err != nil {
		return nil, err
	}
	return &AmendOrderResult{AmendOrderResult: *res}, nil
}

// EditOrder cancel-replaces an order, producing a new tx id.
func (c *Client) EditOrder(pair, txid string, params map[string]string) (*EditOrderResult, error) {
	service := c.NewKrakenClient().NewEditOrderService(pair, txid)
	if v := params["volume"]; v != "" {
		p, err := decimal.NewFromString(v)
		if err != nil {
			return nil, err
		}
		service.SetVolume(p)
	}
	if v := params["price"]; v != "" {
		p, err := decimal.NewFromString(v)
		if err != nil {
			return nil, err
		}
		service.SetPrice(p)
	}
	if v := params["price2"]; v != "" {
		p, err := decimal.NewFromString(v)
		if err != nil {
			return nil, err
		}
		service.SetPrice2(p)
	}
	if v := params["oflags"]; v != "" {
		service.SetOrderFlags(v)
	}
	res, err := service.Do(ctx())
	if err != nil {
		return nil, err
	}
	return &EditOrderResult{EditOrderResult: *res}, nil
}

// CancelOrder cancels a single order by tx id or user reference.
func (c *Client) CancelOrder(txid, clientOrderID string) (*CancelResult, error) {
	service := c.NewKrakenClient().NewCancelOrderService(txid)
	if clientOrderID != "" {
		service.SetClientOrderID(clientOrderID)
	}
	res, err := service.Do(ctx())
	if err != nil {
		return nil, err
	}
	return &CancelResult{Count: res.Count, Pending: res.Pending}, nil
}

// CancelAll cancels all open orders.
func (c *Client) CancelAll() (*CancelResult, error) {
	res, err := c.NewKrakenClient().NewCancelAllOrdersService().Do(ctx())
	if err != nil {
		return nil, err
	}
	return &CancelResult{Count: res.Count}, nil
}

// CancelAllAfter arms (or disarms, with timeout 0) the dead man's switch that
// cancels all orders after the given number of seconds.
func (c *Client) CancelAllAfter(timeout int) (*CancelAllAfterResult, error) {
	res, err := c.NewKrakenClient().NewCancelAllOrdersAfterService(timeout).Do(ctx())
	if err != nil {
		return nil, err
	}
	return &CancelAllAfterResult{CancelAllOrdersAfterResult: *res}, nil
}

// CancelBatch cancels multiple orders by tx id or user reference (max 50).
func (c *Client) CancelBatch(txids []string) (*CancelResult, error) {
	res, err := c.NewKrakenClient().NewCancelOrderBatchService(txids...).Do(ctx())
	if err != nil {
		return nil, err
	}
	return &CancelResult{Count: res.Count}, nil
}
