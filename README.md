# kraken-cli

A command-line interface for the [Kraken](https://docs.kraken.com/api/) **Spot**
REST API, built on top of [`go-kraken`](https://github.com/UnipayFI/go-kraken).

Phase 1 covers the main **private** (API-key authenticated) spot endpoints:
account data, trading, and funding. Public market-data endpoints are out of
scope.

The tool is self-documenting: every subcommand supports `-h/--help` and includes
a link to the corresponding official Kraken API doc.

## Install

```bash
go build -o kraken-cli .
```

## Configuration

Credentials are read from the environment (from the Kraken
[API-management page](https://www.kraken.com/u/security/api)):

```bash
export KRAKEN_API_KEY=...      # API key
export KRAKEN_API_SECRET=...   # base64-encoded private key
export KRAKEN_BASE_URL=...     # optional: override the REST base URL
```

To route traffic through a proxy, use the standard proxy environment variables
(`HTTPS_PROXY`, `HTTP_PROXY`, `NO_PROXY`; lowercase variants are also honoured),
which the HTTP transport reads automatically:

```bash
export HTTPS_PROXY=http://user:pass@host:port   # http/https/socks5 supported
```

## Output

By default results are rendered as a table. Pass the global `--json` flag on any
command to emit the raw response as indented JSON instead:

```bash
kraken-cli account balance --json
```

## Commands

```
kraken-cli
├── account            Account data
│   ├── balance            Cash balances (non-zero)
│   ├── balance-ex         Balances with amount held for orders
│   ├── trade-balance      Margin / equity summary
│   ├── open-orders        Currently open orders
│   ├── closed-orders      Closed orders (paginated)
│   ├── query-orders       Specific orders by tx id (max 50)
│   ├── trades-history     Trade / fill history (paginated)
│   ├── query-trades       Specific trades by tx id (max 20)
│   ├── open-positions     Open margin positions
│   ├── ledgers            Ledger entries (paginated)
│   ├── query-ledgers      Specific ledger entries by id (max 20)
│   ├── trade-volume       30-day volume and fee tier
│   └── api-key-info       Metadata about the signing API key
│
├── order              Trading
│   ├── add                Place a new order (use --validate to dry-run)
│   ├── amend              Amend an open order in place
│   ├── edit               Cancel-replace an order (new tx id)
│   ├── cancel             Cancel a single order
│   ├── cancel-all         Cancel all open orders
│   ├── cancel-all-after   Arm/disarm the dead man's switch
│   └── cancel-batch       Cancel multiple orders by tx id (max 50)
│
└── funding            Deposits, withdrawals, transfers
    ├── deposit-methods    Deposit methods for an asset
    ├── deposit-addresses  Deposit addresses (--new to generate)
    ├── deposit-status     Status of recent deposits
    ├── withdraw-methods   Withdrawal methods
    ├── withdraw-addresses Saved withdrawal addresses
    ├── withdraw-info      Preview a withdrawal's fee and limit
    ├── withdraw           Withdraw to a saved key (irreversible)
    ├── withdraw-status    Status of recent withdrawals
    ├── withdraw-cancel    Cancel a recent withdrawal
    └── wallet-transfer    Transfer between Spot and Futures wallets
```

## Examples

```bash
# Show non-zero balances
kraken-cli account balance

# Validate an order without submitting it
kraken-cli order add --pair XBTUSDT --side buy --type limit \
  --volume 0.0001 --price 30000 --oflags post --validate

# Place, list, then cancel a post-only limit order
kraken-cli order add --pair XBTUSDT --side buy --type limit \
  --volume 0.0001 --price 20000 --oflags post
kraken-cli account open-orders
kraken-cli order cancel-all

# Preview a withdrawal
kraken-cli funding withdraw-info --asset USDT --key MyKey --amount 10
```
