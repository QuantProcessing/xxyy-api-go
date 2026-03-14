# xxyy-api-go

English | [中文](./README_ZH.md)

Go SDK for the [XXYY](https://www.xxyy.io) Open API — on-chain token trading and data queries on Solana, Ethereum, BSC, and Base.

> ⚠️ **API Key = Wallet Access** — Your XXYY API Key can execute real on-chain trades using your wallet balance. Never share it, never commit it to git, never paste it in public channels. If you suspect a leak, regenerate the key immediately at [xxyy.io](https://www.xxyy.io).

## Install

### Claude Code Plugin

```
/plugin install https://github.com/QuantProcessing/xxyy-api-go
```

Or via marketplace:
```
/plugin marketplace add QuantProcessing/xxyy-api-go
/plugin install xxyy-trade-go@xxyy-api-go
```

### Go SDK

```bash
go get github.com/QuantProcessing/xxyy-api-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "os"

    xxyy "github.com/QuantProcessing/xxyy-api-go"
)

func main() {
    client := xxyy.NewClient(os.Getenv("XXYY_API_KEY"))

    // Verify API key
    if err := client.Ping(context.Background()); err != nil {
        panic(err)
    }
    fmt.Println("pong — API Key is valid.")
}
```

## Configuration

```go
// Set your API key as an environment variable
export XXYY_API_KEY=xxyy_ak_your_key_here

// Optional: custom base URL
export XXYY_API_BASE_URL=https://www.xxyy.io
```

### Client Options

```go
client := xxyy.NewClient(apiKey,
    xxyy.WithBaseURL("https://custom.api.url"),
    xxyy.WithTimeout(60 * time.Second),
    xxyy.WithHTTPClient(customHTTPClient),
)
```

## API Reference

### Ping

```go
err := client.Ping(ctx)
```

### Buy Token

```go
slippage := 20
resp, err := client.BuyToken(ctx, &xxyy.SwapRequest{
    Chain:         xxyy.ChainSOL,
    WalletAddress: "your_wallet_address",
    TokenAddress:  "token_contract_address",
    Amount:        0.1,    // Amount in native currency (SOL/ETH/BNB)
    Tip:           0.001,  // SOL: 0.001-0.1 SOL; EVM: 0.1-100 Gwei
    Slippage:      &slippage,
})
fmt.Printf("TxID: %s\n", resp.TxID)
fmt.Printf("Explorer: %s\n", xxyy.ExplorerURL(xxyy.ChainSOL, resp.TxID))
```

### Sell Token

```go
resp, err := client.SellToken(ctx, &xxyy.SwapRequest{
    Chain:         xxyy.ChainSOL,
    WalletAddress: "your_wallet_address",
    TokenAddress:  "token_contract_address",
    Amount:        50,     // Sell percentage (1-100)
    Tip:           0.001,
})
```

### Query Trade

```go
trade, err := client.QueryTrade(ctx, "transaction_id")
fmt.Printf("Status: %s\n", trade.Status) // pending / success / failed
```

### Feed Scan

Scan Meme token lists. Only supports **SOL** and **BSC** chains.

```go
data, err := client.FeedScan(ctx, xxyy.FeedNew, xxyy.ChainSOL, &xxyy.FeedFilter{
    MC:     "10000,500000",  // market cap range
    Holder: "50,",           // min 50 holders
})
for _, item := range data.Items {
    fmt.Printf("%s ($%.6f) — MCap: $%.0f\n", item.Symbol, item.PriceUSD, item.MarketCapUSD)
}
```

Feed types: `FeedNew`, `FeedAlmost`, `FeedCompleted`

### Token Query

```go
token, err := client.TokenQuery(ctx, "contract_address", xxyy.ChainSOL)
fmt.Printf("Price: $%v\n", token.TradeInfo.Price)
fmt.Printf("HoneyPot: %v\n", token.SecurityInfo.HoneyPot)
```

### List Wallets

```go
wallets, err := client.ListWallets(ctx, &xxyy.WalletsRequest{
    Chain: xxyy.ChainSOL,
})
for _, w := range wallets.List {
    fmt.Printf("%s: %v %s\n", w.Name, w.Balance, xxyy.NativeToken(xxyy.ChainSOL))
}
```

### Wallet Info

```go
info, err := client.WalletInfo(ctx, &xxyy.WalletInfoRequest{
    WalletAddress: "your_wallet_address",
    Chain:         xxyy.ChainSOL,
    TokenAddress:  "optional_token_address", // optional: check token holdings
})
fmt.Printf("Balance: %v\n", info.Balance)
```

## Supported Chains

| Chain | Constant | Native Token |
|-------|----------|-------------|
| Solana | `xxyy.ChainSOL` | SOL |
| Ethereum | `xxyy.ChainETH` | ETH |
| BSC | `xxyy.ChainBSC` | BNB |
| Base | `xxyy.ChainBase` | ETH |

## Error Handling

```go
resp, err := client.BuyToken(ctx, req)
if err != nil {
    var xxErr *xxyy.XxyyError
    if errors.As(err, &xxErr) {
        if xxErr.IsAPIKeyError() {
            // API key invalid or disabled
        }
        if xxErr.IsRateLimited() {
            // Rate limited (auto-retry exhausted)
        }
        if xxErr.IsServerError() {
            // Server-side error
        }
    }
}
```

## Examples

See the [examples/](./examples/) directory for runnable examples:

- [ping](./examples/ping/) — verify API key
- [buy](./examples/buy/) — buy a token and query trade status
- [feed](./examples/feed/) — scan tokens with filters
- [wallets](./examples/wallets/) — list wallets

## License

[MIT](./LICENSE)
