package funding

import (
	"context"

	"github.com/UnipayFI/kraken-cli/exchange"
	"github.com/shopspring/decimal"
)

// Client wraps the shared exchange client with funding methods.
type Client struct {
	*exchange.Client
}

func ctx() context.Context { return context.Background() }

// GetDepositMethods returns the deposit methods available for an asset.
func (c *Client) GetDepositMethods(asset string) (*DepositMethodList, error) {
	methods, err := c.NewKrakenClient().NewGetDepositMethodsService(asset).Do(ctx())
	if err != nil {
		return nil, err
	}
	list := DepositMethodList(methods)
	return &list, nil
}

// GetDepositAddresses returns (or generates) deposit addresses for an asset and
// method.
func (c *Client) GetDepositAddresses(asset, method string, newAddr bool) (*DepositAddressList, error) {
	service := c.NewKrakenClient().NewGetDepositAddressesService(asset, method)
	if newAddr {
		service.SetNew(true)
	}
	addrs, err := service.Do(ctx())
	if err != nil {
		return nil, err
	}
	list := DepositAddressList(addrs)
	return &list, nil
}

// GetDepositStatus returns the status of recent deposits.
func (c *Client) GetDepositStatus(asset, method string) (*TransferStatusList, error) {
	service := c.NewKrakenClient().NewGetDepositStatusService()
	if asset != "" {
		service.SetAsset(asset)
	}
	if method != "" {
		service.SetMethod(method)
	}
	statuses, err := service.Do(ctx())
	if err != nil {
		return nil, err
	}
	list := TransferStatusList(statuses)
	return &list, nil
}

// GetWithdrawMethods returns the withdrawal methods available for an asset.
func (c *Client) GetWithdrawMethods(asset, network string) (*WithdrawMethodList, error) {
	service := c.NewKrakenClient().NewGetWithdrawMethodsService()
	if asset != "" {
		service.SetAsset(asset)
	}
	if network != "" {
		service.SetNetwork(network)
	}
	methods, err := service.Do(ctx())
	if err != nil {
		return nil, err
	}
	list := WithdrawMethodList(methods)
	return &list, nil
}

// GetWithdrawAddresses returns the saved withdrawal addresses.
func (c *Client) GetWithdrawAddresses(asset, method string) (*WithdrawAddressList, error) {
	service := c.NewKrakenClient().NewGetWithdrawAddressesService()
	if asset != "" {
		service.SetAsset(asset)
	}
	if method != "" {
		service.SetMethod(method)
	}
	addrs, err := service.Do(ctx())
	if err != nil {
		return nil, err
	}
	list := WithdrawAddressList(addrs)
	return &list, nil
}

// GetWithdrawInfo previews the fee and limit for a withdrawal without executing
// it.
func (c *Client) GetWithdrawInfo(asset, key string, amount decimal.Decimal) (*WithdrawInfo, error) {
	info, err := c.NewKrakenClient().NewGetWithdrawInfoService(asset, key, amount).Do(ctx())
	if err != nil {
		return nil, err
	}
	return &WithdrawInfo{WithdrawInfo: *info}, nil
}

// Withdraw executes a withdrawal to a pre-configured withdrawal key.
func (c *Client) Withdraw(asset, key string, amount decimal.Decimal, address string) (*RefResult, error) {
	service := c.NewKrakenClient().NewWithdrawService(asset, key, amount)
	if address != "" {
		service.SetAddress(address)
	}
	ref, err := service.Do(ctx())
	if err != nil {
		return nil, err
	}
	return &RefResult{RefID: ref.RefID}, nil
}

// GetWithdrawStatus returns the status of recent withdrawals.
func (c *Client) GetWithdrawStatus(asset, method string) (*TransferStatusList, error) {
	service := c.NewKrakenClient().NewGetWithdrawStatusService()
	if asset != "" {
		service.SetAsset(asset)
	}
	if method != "" {
		service.SetMethod(method)
	}
	statuses, err := service.Do(ctx())
	if err != nil {
		return nil, err
	}
	list := TransferStatusList(statuses)
	return &list, nil
}

// WithdrawCancel cancels a recently requested withdrawal by reference id.
func (c *Client) WithdrawCancel(asset, refID string) (bool, error) {
	return c.NewKrakenClient().NewWithdrawCancelService(asset, refID).Do(ctx())
}

// WalletTransfer moves funds between the Spot and Futures wallets.
func (c *Client) WalletTransfer(asset string, amount decimal.Decimal, from, to string) (*RefResult, error) {
	service := c.NewKrakenClient().NewWalletTransferService(asset, amount)
	if from != "" {
		service.SetFrom(from)
	}
	if to != "" {
		service.SetTo(to)
	}
	ref, err := service.Do(ctx())
	if err != nil {
		return nil, err
	}
	return &RefResult{RefID: ref.RefID}, nil
}
