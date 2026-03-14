---
name: xxyy-trade-go
description: >-
  This skill should be used when the user asks to interact with the XXYY
  trading platform using Go code. It covers buying/selling tokens, querying
  trades, scanning token feeds, querying token details, and managing wallets
  on Solana/ETH/BSC/Base chains via the xxyy-api-go SDK.
  Triggers: "buy token go", "sell token go", "xxyy go sdk", "trade crypto go",
  "scan tokens go", "query token go", "list wallets go", "wallet balance go".
version: 0.1.1
---

# XXYY Trade Go SDK

On-chain token trading and data queries on Solana, Ethereum, BSC, and Base via the [xxyy-api-go](https://github.com/QuantProcessing/xxyy-api-go) SDK.

## Prerequisites

1. Go 1.21+
2. Install the SDK:
```bash
go get github.com/QuantProcessing/xxyy-api-go@latest
```
3. Set environment variable:
   - `XXYY_API_KEY` (required) — format: `xxyy_ak_xxxx`. Get one at https://www.xxyy.io/apikey

## Security Notes

- **⚠️ API Key = Wallet access** — Your XXYY API Key can execute real on-chain trades using your wallet balance. If it leaks, anyone can buy/sell tokens with your funds. Never share it, never commit it to version control, never expose it in logs.
- **Custodial trading model** — XXYY executes trades on your behalf. No private keys or wallet signing needed.
- **No read-only mode** — The same API Key is used for both data queries and trading.

## SDK Import & Client Setup

```go
import xxyy "github.com/QuantProcessing/xxyy-api-go"

// Create client
client := xxyy.NewClient(os.Getenv("XXYY_API_KEY"))

// With options
client := xxyy.NewClient(apiKey,
    xxyy.WithBaseURL("https://custom.url"),
    xxyy.WithTimeout(60 * time.Second),
    xxyy.WithHTTPClient(customHTTPClient),
)
```

All methods accept `context.Context` as the first argument.

## API Reference

> **STRICT: Only the methods listed below exist. Do NOT guess or construct any API call not documented here.**

### Buy Token

```go
slippage := 20
resp, err := client.BuyToken(ctx, &xxyy.SwapRequest{
    Chain:         xxyy.ChainSOL,       // sol / eth / bsc / base
    WalletAddress: "<user_wallet>",
    TokenAddress:  "<token_contract>",
    Amount:        0.1,                 // Amount in native currency (SOL/ETH/BNB)
    Tip:           0.001,               // Priority fee
    Slippage:      &slippage,           // Optional, default 20
})
// resp.TxID contains the transaction ID
```

#### Buy Parameters

| Param | Required | Type | Valid values | Description |
|-------|----------|------|-------------|-------------|
| `Chain` | YES | `xxyy.Chain` | `ChainSOL` / `ChainETH` / `ChainBSC` / `ChainBase` | Only these 4 values |
| `WalletAddress` | YES | string | SOL: Base58 32-44 chars; EVM: 0x+40hex | Must match chain |
| `TokenAddress` | YES | string | Valid contract address | Token to buy |
| `Amount` | YES | float64 | > 0 | Amount in native currency |
| `Tip` | YES | float64 | SOL: 0.001-0.1 (SOL); EVM: 0.1-100 (Gwei) | Priority fee |
| `Slippage` | NO | *int | 0-100 | Default 20 |
| `Model` | NO | *int | 1 or 2 | 1=anti-sandwich (default), 2=fast |
| `PriorityFee` | NO | *float64 | >= 0 | Solana only |

### Sell Token

```go
resp, err := client.SellToken(ctx, &xxyy.SwapRequest{
    Chain:         xxyy.ChainSOL,
    WalletAddress: "<user_wallet>",
    TokenAddress:  "<token_contract>",
    Amount:        50,                  // Sell percentage: 50 = sell 50%
    Tip:           0.001,
})
```

Sell `Amount` is a **percentage** (1-100), not an absolute value.

### tip / priorityFee Rules

- `Tip` (required) — Universal priority fee for ALL chains.
  - SOL: unit is SOL (0.001 - 0.1 recommended)
  - EVM (eth/bsc/base): unit is Gwei (0.1 - 100 recommended)
- `PriorityFee` (optional) — Only effective on Solana. Extra fee in addition to Tip.
- **SDK validates tip range** — returns error if outside recommended range.

### Query Trade

```go
trade, err := client.QueryTrade(ctx, "<txId>")
// trade.Status: int (1=pending, 2=success, 3=failed)
// trade.StatusText(): "pending" / "success" / "failed"
// trade.IsSuccess(), trade.IsPending(), trade.IsFailed()
// trade.IsBuy: bool
// trade.BaseAmount, trade.QuoteAmount: float64
```

### Ping

```go
err := client.Ping(ctx)
// nil = API Key valid
```

### Feed (Scan Tokens)

```go
data, err := client.FeedScan(ctx, xxyy.FeedNew, xxyy.ChainSOL, &xxyy.FeedFilter{
    MC:     "10000,500000",  // Market cap range USD
    Holder: "50,",           // Min 50 holders
})
for _, item := range data.Items {
    // item.Symbol, item.PriceUSD, item.MarketCapUSD, item.Holders
    // item.LaunchPlatform.Name, item.LaunchPlatform.Progress
    // item.DevHoldPercent, item.Snipers, item.HasLink
}
```

#### Feed Parameters

| Param | Required | Type | Valid values | Description |
|-------|----------|------|-------------|-------------|
| feedType | YES | `xxyy.FeedType` | `FeedNew` / `FeedAlmost` / `FeedCompleted` | Token list type |
| chain | YES | `xxyy.Chain` | `ChainSOL` / `ChainBSC` only | **Only sol and bsc** |
| filter | NO | `*xxyy.FeedFilter` | nil for no filters | See FeedFilter struct |

#### FeedFilter Fields (all optional)

**Array filters**: `Dex`, `QuoteTokens`, `Link`, `Keywords`, `IgnoreWords`

**Range filters** ("min,max" format): `MC`, `Liq`, `Vol`, `Holder`, `CreateTime`, `TradeCount`, `BuyCount`, `SellCount`, `DevBuy`, `DevSell`, `DevHp`, `TopHp`, `InsiderHp`, `BundleHp`, `NewWalletHp`, `Progress`, `Snipers`, `XnameCount`, `TagHolder`, `KOL`

**Integer flags**: `DexPay`, `OneLink`, `Live` (1=enabled)

#### DEX Values by Chain
- **SOL**: `pump`, `pumpmayhem`, `bonk`, `heaven`, `believe`, `daosfun`, `launchlab`, `mdbc`, `jupstudio`, `mdbcbags`, `trends`, `moonshotn`, `boop`, `moon`, `time`
- **BSC**: `four`, `four_agent`, `bnonly`, `flap`

### Token Query

```go
token, err := client.TokenQuery(ctx, "<contract_address>", xxyy.ChainSOL)
// token.BaseSymbol, token.ChainID
// token.TradeInfo.Price, .MarketCapUSD, .Holder, .HourTradeNum, .HourTradeVolume
// token.SecurityInfo.HoneyPot, .OpenSource, .NoOwner, .Locked
// token.TaxInfo.Buy, .Sell  (percentage strings)
// token.LinkInfo.TG, .X, .Web
// token.Dev.Address, .Pct
// token.TopHolderPct, token.TopHolderList
```

Supports all 4 chains. Default chain is `ChainSOL`.

### List Wallets

```go
wallets, err := client.ListWallets(ctx, &xxyy.WalletsRequest{
    Chain:        xxyy.ChainSOL,
    PageNum:      1,                    // Optional
    PageSize:     20,                   // Optional, max 20
    TokenAddress: "<token_contract>",   // Optional: check token holdings
})
for _, w := range wallets.List {
    // w.Name, w.Address, w.Balance
    // w.IsPinned() — true if wallet is pinned (⭐)
    // w.IsImport — true if imported wallet
    // w.TokenBalance — only present when TokenAddress is set
}
```

### Wallet Info

```go
info, err := client.WalletInfo(ctx, &xxyy.WalletInfoRequest{
    WalletAddress: "<wallet_address>",  // Required
    Chain:         xxyy.ChainSOL,       // Optional, default sol
    TokenAddress:  "<token_contract>",  // Optional: check token holdings
})
// info.Address, info.Name, info.Balance
// info.TokenBalance.UIAmount, .Amount, .Decimals
```

## Chain Constants

| Chain | Constant | Native Token |
|-------|----------|-------------|
| Solana | `xxyy.ChainSOL` | SOL |
| Ethereum | `xxyy.ChainETH` | ETH |
| BSC | `xxyy.ChainBSC` | BNB |
| Base | `xxyy.ChainBase` | ETH |

**Helper functions:**
- `xxyy.ExplorerURL(chain, txID)` — returns block explorer URL
- `xxyy.NativeToken(chain)` — returns native token symbol ("SOL", "ETH", "BNB")
- `chain.IsValid()` — validates chain value
- `chain.IsFeedSupported()` — true for SOL/BSC only

## Error Handling

```go
import "errors"

resp, err := client.BuyToken(ctx, req)
if err != nil {
    var xxErr *xxyy.XxyyError
    if errors.As(err, &xxErr) {
        switch {
        case xxErr.IsAPIKeyError():   // 8060/8061: invalid/disabled key
        case xxErr.IsRateLimited():   // 8062: rate limited
        case xxErr.IsServerError():   // 300: server error
        }
        fmt.Printf("Code: %d, Message: %s\n", xxErr.Code, xxErr.Message)
    }
}
```

| Code | Meaning | SDK Behavior |
|------|---------|-------------|
| 8060 | API Key invalid | Return error immediately |
| 8061 | API Key disabled | Return error immediately |
| 8062 | Rate limited | Auto-retry 2x with 2s delay (swap: no retry) |
| 300 | Server error | Return error immediately |

## Execution Rules

1. **Always confirm before trading** — Ask user to confirm: chain, token address, amount/percentage, buy or sell
2. **Auto-query wallet** — If user does not provide wallet address:
   a. Call `ListWallets`. If only 1 wallet, auto-select. If multiple, ask user. If none, guide to https://www.xxyy.io/wallet/manager?chainId={chain}
   b. Remember selected wallet as default for that chain in session
3. **Poll trade result** — After swap, call `QueryTrade` up to 3 times with 5s intervals
4. **Show transaction link** — Always display `xxyy.ExplorerURL(chain, txID)`
5. **Never retry failed swaps** — Show error to user instead
6. **Chain-wallet validation** — SDK validates wallet address format automatically
7. **Strict parameter validation** — SDK validates all parameters before sending:
   - Chain must be valid
   - Address format (Base58 for SOL, 0x for EVM)
   - Tip within recommended range
   - Buy amount > 0; sell amount 1-100
   - Model must be 1 or 2
   - PriorityFee only on SOL

## Feed Rules

1. **Feed only supports `ChainSOL` and `ChainBSC`** — SDK returns error for other chains
2. **Feed types**: `FeedNew` (new), `FeedAlmost` (near graduation), `FeedCompleted` (graduated)
3. **No auto-trading** — Feed is for observation only. NEVER automatically buy based on scan results

## Token Query Rules

1. **HoneyPot warning** — If `token.SecurityInfo.HoneyPot == true`, display prominent warning
2. **High tax alert** — If buy/sell tax > 5%, warn user
3. **Display format** — Present in groups: Trade Info → Security → Tax → Holders → Links

## Wallet Address Formats

| Chain | Format | Regex |
|-------|--------|-------|
| SOL | Base58, 32-44 chars | `^[1-9A-HJ-NP-Za-km-z]{32,44}$` |
| ETH/BSC/Base | 0x + 40 hex | `^0x[0-9a-fA-F]{40}$` |

## Block Explorer URLs

- SOL: `https://solscan.io/tx/{txId}`
- ETH: `https://etherscan.io/tx/{txId}`
- BSC: `https://bscscan.com/tx/{txId}`
- BASE: `https://basescan.org/tx/{txId}`

Or use: `xxyy.ExplorerURL(chain, txID)`

## Default Trade Parameters

| Chain | Slippage | Tip |
|-------|----------|-----|
| SOL | 20% | 0.001 SOL |
| ETH/BSC/Base | 20% | 1 Gwei |

## Complete Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"

    xxyy "github.com/QuantProcessing/xxyy-api-go"
)

func main() {
    client := xxyy.NewClient(os.Getenv("XXYY_API_KEY"))
    ctx := context.Background()

    // 1. Verify API key
    if err := client.Ping(ctx); err != nil {
        panic(err)
    }

    // 2. List wallets
    wallets, _ := client.ListWallets(ctx, &xxyy.WalletsRequest{Chain: xxyy.ChainBSC})
    wallet := wallets.List[0].Address

    // 3. Query token safety
    token, _ := client.TokenQuery(ctx, "0xTokenAddr", xxyy.ChainBSC)
    if token.SecurityInfo != nil && token.SecurityInfo.HoneyPot {
        fmt.Println("⛔ HONEYPOT! Aborting.")
        return
    }

    // 4. Buy token
    slippage := 20
    buy, _ := client.BuyToken(ctx, &xxyy.SwapRequest{
        Chain: xxyy.ChainBSC, WalletAddress: wallet,
        TokenAddress: "0xTokenAddr", Amount: 0.01,
        Tip: 1, Slippage: &slippage,
    })
    fmt.Printf("Buy TX: %s\n", xxyy.ExplorerURL(xxyy.ChainBSC, buy.TxID))

    // 5. Poll status
    for i := 0; i < 3; i++ {
        time.Sleep(5 * time.Second)
        trade, _ := client.QueryTrade(ctx, buy.TxID)
        fmt.Printf("Status: %s\n", trade.StatusText())
        if trade.IsSuccess() || trade.IsFailed() {
            break
        }
    }

    // 6. Sell 100%
    sell, _ := client.SellToken(ctx, &xxyy.SwapRequest{
        Chain: xxyy.ChainBSC, WalletAddress: wallet,
        TokenAddress: "0xTokenAddr", Amount: 100, Tip: 1,
    })
    fmt.Printf("Sell TX: %s\n", xxyy.ExplorerURL(xxyy.ChainBSC, sell.TxID))
}
```
