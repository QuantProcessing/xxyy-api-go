# xxyy-api-go

[English](./README.md) | 中文

XXYY Open API 的 Go SDK — 支持在 Solana、Ethereum、BSC、Base 链上进行链上代币交易和数据查询。

> ⚠️ **API Key = 钱包权限** — 你的 XXYY API Key 可以使用你的钱包余额执行真实的链上交易。请勿分享、提交到 git 或粘贴到公开渠道。如果怀疑泄露，请立即在 [xxyy.io](https://www.xxyy.io) 重新生成密钥。

## 安装

```bash
go get github.com/QuantProcessing/xxyy-api-go
```

## 快速开始

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

    // 验证 API Key
    if err := client.Ping(context.Background()); err != nil {
        panic(err)
    }
    fmt.Println("pong — API Key 有效")
}
```

## 配置

```bash
# 设置 API Key 环境变量
export XXYY_API_KEY=xxyy_ak_your_key_here

# 可选：自定义 API 地址
export XXYY_API_BASE_URL=https://www.xxyy.io
```

### 客户端选项

```go
client := xxyy.NewClient(apiKey,
    xxyy.WithBaseURL("https://custom.api.url"),       // 自定义 API 地址
    xxyy.WithTimeout(60 * time.Second),                // 自定义超时（默认 30s）
    xxyy.WithHTTPClient(customHTTPClient),             // 自定义 HTTP 客户端
)
```

## API 说明

### Ping（健康检查）

```go
err := client.Ping(ctx)
```

### 买入代币

```go
slippage := 20
resp, err := client.BuyToken(ctx, &xxyy.SwapRequest{
    Chain:         xxyy.ChainSOL,
    WalletAddress: "你的钱包地址",
    TokenAddress:  "代币合约地址",
    Amount:        0.1,    // 原生代币数量 (SOL/ETH/BNB)
    Tip:           0.001,  // SOL: 0.001-0.1 SOL; EVM: 0.1-100 Gwei
    Slippage:      &slippage,
})
fmt.Printf("交易ID: %s\n", resp.TxID)
fmt.Printf("浏览器: %s\n", xxyy.ExplorerURL(xxyy.ChainSOL, resp.TxID))
```

### 卖出代币

```go
resp, err := client.SellToken(ctx, &xxyy.SwapRequest{
    Chain:         xxyy.ChainSOL,
    WalletAddress: "你的钱包地址",
    TokenAddress:  "代币合约地址",
    Amount:        50,     // 卖出百分比 (1-100)
    Tip:           0.001,
})
```

### 查询交易状态

```go
trade, err := client.QueryTrade(ctx, "交易ID")
fmt.Printf("状态: %s\n", trade.Status) // pending / success / failed
```

### 扫描代币 (Feed)

扫描 Meme 代币列表。仅支持 **SOL** 和 **BSC** 链。

```go
data, err := client.FeedScan(ctx, xxyy.FeedNew, xxyy.ChainSOL, &xxyy.FeedFilter{
    MC:     "10000,500000",  // 市值范围
    Holder: "50,",           // 最少 50 个持有人
})
for _, item := range data.Items {
    fmt.Printf("%s ($%.6f) — 市值: $%.0f\n", item.Symbol, item.PriceUSD, item.MarketCapUSD)
}
```

扫描类型：`FeedNew`（新发行）、`FeedAlmost`（将毕业）、`FeedCompleted`（已毕业）

### 代币详情查询

```go
token, err := client.TokenQuery(ctx, "合约地址", xxyy.ChainSOL)
fmt.Printf("价格: $%v\n", token.TradeInfo.Price)
fmt.Printf("蜜罐: %v\n", token.SecurityInfo.HoneyPot)
```

### 钱包列表

```go
wallets, err := client.ListWallets(ctx, &xxyy.WalletsRequest{
    Chain: xxyy.ChainSOL,
})
for _, w := range wallets.List {
    fmt.Printf("%s: %v %s\n", w.Name, w.Balance, xxyy.NativeToken(xxyy.ChainSOL))
}
```

### 钱包详情

```go
info, err := client.WalletInfo(ctx, &xxyy.WalletInfoRequest{
    WalletAddress: "你的钱包地址",
    Chain:         xxyy.ChainSOL,
    TokenAddress:  "代币地址", // 可选：查看代币持仓
})
fmt.Printf("余额: %v\n", info.Balance)
```

## 支持的链

| 链 | 常量 | 原生代币 |
|----|------|----------|
| Solana | `xxyy.ChainSOL` | SOL |
| Ethereum | `xxyy.ChainETH` | ETH |
| BSC | `xxyy.ChainBSC` | BNB |
| Base | `xxyy.ChainBase` | ETH |

## 错误处理

```go
resp, err := client.BuyToken(ctx, req)
if err != nil {
    var xxErr *xxyy.XxyyError
    if errors.As(err, &xxErr) {
        if xxErr.IsAPIKeyError() {
            // API Key 无效或已禁用
        }
        if xxErr.IsRateLimited() {
            // 限流（自动重试已用尽）
        }
        if xxErr.IsServerError() {
            // 服务端错误
        }
    }
}
```

| 错误码 | 含义 | 处理方式 |
|--------|------|----------|
| 8060 | API Key 无效 | 到 xxyy.io 重新生成 |
| 8061 | API Key 已禁用 | 到 xxyy.io 重新生成 |
| 8062 | 限流 | 自动重试 2 次（swap 除外） |
| 300 | 服务器错误 | 稍后重试 |

## 示例

查看 [examples/](./examples/) 目录获取可运行的示例：

- [ping](./examples/ping/) — 验证 API Key
- [buy](./examples/buy/) — 买入代币并查询交易状态
- [feed](./examples/feed/) — 使用过滤器扫描代币
- [wallets](./examples/wallets/) — 列出钱包

## 许可证

[MIT](./LICENSE)
