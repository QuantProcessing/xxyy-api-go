package xxyy

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ApiResponse is the unified response wrapper for all XXYY API responses.
type ApiResponse[T any] struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    T      `json:"data"`
	Success bool   `json:"success"`
}

// ---------------------------------------------------------------------------
// Swap (Buy / Sell)
// ---------------------------------------------------------------------------

// SwapRequest represents the request body for buying or selling a token.
type SwapRequest struct {
	Chain         Chain   `json:"chain"`
	WalletAddress string  `json:"walletAddress"`
	TokenAddress  string  `json:"tokenAddress"`
	Amount        float64 `json:"amount"`
	Tip           float64 `json:"tip"`
	Slippage      *int    `json:"slippage,omitempty"`
	Model         *int    `json:"model,omitempty"`
	PriorityFee   *float64 `json:"priorityFee,omitempty"`
}

// swapBody is the internal request body sent to the API.
type swapBody struct {
	Chain         Chain   `json:"chain"`
	WalletAddress string  `json:"walletAddress"`
	TokenAddress  string  `json:"tokenAddress"`
	IsBuy         bool    `json:"isBuy"`
	Amount        float64 `json:"amount"`
	Tip           float64 `json:"tip"`
	Slippage      *int    `json:"slippage,omitempty"`
	Model         *int    `json:"model,omitempty"`
	PriorityFee   *float64 `json:"priorityFee,omitempty"`
}

// SwapResponse represents the response from a swap (buy/sell) operation.
type SwapResponse struct {
	TxID      string `json:"txId,omitempty"`
	Signature string `json:"signature,omitempty"`
}

// ---------------------------------------------------------------------------
// Trade Query
// ---------------------------------------------------------------------------

// Trade status constants (numeric values from API).
// Based on live API observation: 1=pending, 2=success, 3=failed.
const (
	TradeStatusPending = 1
	TradeStatusSuccess = 2
	TradeStatusFailed  = 3
)

var tradeStatusText = map[int]string{
	TradeStatusPending: "pending",
	TradeStatusSuccess: "success",
	TradeStatusFailed:  "failed",
}

// TradeData represents the response from querying a trade status.
// Note: The XXYY API returns status as a number (0=pending, 1=failed, 2=success)
// and isBuy as a number (0=sell, 1=buy), despite documentation suggesting strings/booleans.
type TradeData struct {
	TxID          string  `json:"txId"`
	Status        int     `json:"-"`               // Parsed from number or string
	StatusDesc    string  `json:"statusDesc,omitempty"`
	Chain         string  `json:"chain,omitempty"`
	TokenAddress  string  `json:"tokenAddress,omitempty"`
	WalletAddress string  `json:"walletAddress,omitempty"`
	IsBuy         bool    `json:"-"`               // Parsed from number or boolean
	BaseAmount    float64 `json:"baseAmount,omitempty"`
	QuoteAmount   float64 `json:"quoteAmount,omitempty"`
	CreateTime    string  `json:"createTime,omitempty"`
	UpdateTime    string  `json:"updateTime,omitempty"`
}

// StatusText returns the human-readable status string.
func (t *TradeData) StatusText() string {
	if s, ok := tradeStatusText[t.Status]; ok {
		return s
	}
	return fmt.Sprintf("unknown(%d)", t.Status)
}

// IsSuccess returns true if the trade completed successfully.
func (t *TradeData) IsSuccess() bool {
	return t.Status == TradeStatusSuccess
}

// IsPending returns true if the trade is still pending.
func (t *TradeData) IsPending() bool {
	return t.Status == TradeStatusPending
}

// IsFailed returns true if the trade failed.
func (t *TradeData) IsFailed() bool {
	return t.Status == TradeStatusFailed
}

// UnmarshalJSON implements custom JSON unmarshaling for TradeData.
// Handles the API returning status as number and isBuy as number (0/1).
func (t *TradeData) UnmarshalJSON(data []byte) error {
	type Alias TradeData
	aux := &struct {
		*Alias
		Status json.RawMessage `json:"status"`
		IsBuy  json.RawMessage `json:"isBuy"`
	}{Alias: (*Alias)(t)}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Parse status: can be number (0/1/2) or string ("pending"/"success"/"failed")
	if len(aux.Status) > 0 {
		var statusNum int
		if err := json.Unmarshal(aux.Status, &statusNum); err == nil {
			t.Status = statusNum
		} else {
			var statusStr string
			if err := json.Unmarshal(aux.Status, &statusStr); err == nil {
				switch statusStr {
				case "pending":
					t.Status = TradeStatusPending
				case "failed":
					t.Status = TradeStatusFailed
				case "success":
					t.Status = TradeStatusSuccess
				default:
					n, _ := strconv.Atoi(statusStr)
					t.Status = n
				}
			}
		}
	}

	// Parse isBuy: can be number (0/1) or boolean (true/false)
	if len(aux.IsBuy) > 0 {
		var buyBool bool
		if err := json.Unmarshal(aux.IsBuy, &buyBool); err == nil {
			t.IsBuy = buyBool
		} else {
			var buyNum int
			if err := json.Unmarshal(aux.IsBuy, &buyNum); err == nil {
				t.IsBuy = buyNum == 1
			}
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Feed (Scan Tokens)
// ---------------------------------------------------------------------------

// FeedType represents the type of feed scan.
type FeedType string

const (
	FeedNew       FeedType = "NEW"
	FeedAlmost    FeedType = "ALMOST"
	FeedCompleted FeedType = "COMPLETED"
)

// IsValid returns true if the feed type is valid.
func (f FeedType) IsValid() bool {
	switch f {
	case FeedNew, FeedAlmost, FeedCompleted:
		return true
	}
	return false
}

// FeedFilter represents the optional filter parameters for feed scanning.
// Range parameters use comma-separated string format "min,max".
// Leave one side empty to set only min or max (e.g. "100," = min 100, ",50" = max 50).
type FeedFilter struct {
	Dex         []string `json:"dex,omitempty"`         // DEX platform filter
	QuoteTokens []string `json:"quoteTokens,omitempty"` // Quote token filter
	Link        []string `json:"link,omitempty"`         // Social media link filter
	Keywords    []string `json:"keywords,omitempty"`     // Token name/symbol keyword match
	IgnoreWords []string `json:"ignoreWords,omitempty"`  // Ignore keywords

	// Range parameters ("min,max" format)
	MC          string `json:"mc,omitempty"`          // Market cap range (USD)
	Liq         string `json:"liq,omitempty"`         // Liquidity range (USD)
	Vol         string `json:"vol,omitempty"`         // Trading volume range (USD)
	Holder      string `json:"holder,omitempty"`      // Holder count range
	CreateTime  string `json:"createTime,omitempty"`  // Creation time range (minutes from now)
	TradeCount  string `json:"tradeCount,omitempty"`  // Trade count range
	BuyCount    string `json:"buyCount,omitempty"`    // Buy count range
	SellCount   string `json:"sellCount,omitempty"`   // Sell count range
	DevBuy      string `json:"devBuy,omitempty"`      // Dev buy amount range
	DevSell     string `json:"devSell,omitempty"`     // Dev sell amount range
	DevHp       string `json:"devHp,omitempty"`       // Dev holding % range
	TopHp       string `json:"topHp,omitempty"`       // Top10 holding % range
	InsiderHp   string `json:"insiderHp,omitempty"`   // Insider holding % range
	BundleHp    string `json:"bundleHp,omitempty"`    // Bundle holding % range
	NewWalletHp string `json:"newWalletHp,omitempty"` // New wallet holding % range
	Progress    string `json:"progress,omitempty"`    // Graduation progress % range
	Snipers     string `json:"snipers,omitempty"`     // Sniper count range
	XnameCount  string `json:"xnameCount,omitempty"`  // Twitter rename count range
	TagHolder   string `json:"tagHolder,omitempty"`   // Watched wallet buy count range
	KOL         string `json:"kol,omitempty"`         // KOL buy count range

	// Integer flags
	DexPay  *int `json:"dexPay,omitempty"`  // DexScreener paid, 1=filter paid only
	OneLink *int `json:"oneLink,omitempty"` // At least one social link, 1=enabled
	Live    *int `json:"live,omitempty"`    // Currently live streaming, 1=filter live
}

// FeedItem represents a single token entry in the feed response.
type FeedItem struct {
	TokenAddress    string          `json:"tokenAddress"`
	Symbol          string          `json:"symbol"`
	Name            string          `json:"name"`
	CreateTime      int64           `json:"createTime"`
	DexName         string          `json:"dexName,omitempty"`
	LaunchPlatform  *LaunchPlatform `json:"launchPlatform,omitempty"`
	Holders         int             `json:"holders"`
	PriceUSD        float64         `json:"priceUSD"`
	MarketCapUSD    float64         `json:"marketCapUSD"`
	DevHoldPercent  float64         `json:"devHoldPercent,omitempty"`
	HasLink         bool            `json:"hasLink,omitempty"`
	Snipers         int             `json:"snipers,omitempty"`
	Volume          float64         `json:"volume,omitempty"`
	TradeCount      int             `json:"tradeCount,omitempty"`
	BuyCount        int             `json:"buyCount,omitempty"`
	SellCount       int             `json:"sellCount,omitempty"`
	TopHolderPercent float64        `json:"topHolderPercent,omitempty"`
	InsiderHp       float64         `json:"insiderHp,omitempty"`
	BundleHp        float64         `json:"bundleHp,omitempty"`
	QuoteToken      string          `json:"quoteToken,omitempty"`
}

// LaunchPlatform describes the launch platform details.
type LaunchPlatform struct {
	Name      string `json:"name"`
	Progress  string `json:"progress"`
	Completed bool   `json:"completed"`
}

// FeedData represents the response from a feed scan.
type FeedData struct {
	Items []FeedItem `json:"items"`
}

// ---------------------------------------------------------------------------
// Token Query
// ---------------------------------------------------------------------------

// TokenData represents the detailed information for a token.
type TokenData struct {
	ChainID       string       `json:"chainId"`
	TokenAddress  string       `json:"tokenAddress"`
	BaseSymbol    string       `json:"baseSymbol"`
	TradeInfo     *TradeInfo   `json:"tradeInfo,omitempty"`
	PairInfo      *PairInfo    `json:"pairInfo,omitempty"`
	SecurityInfo  *SecurityInfo `json:"securityInfo,omitempty"`
	TaxInfo       *TaxInfo     `json:"taxInfo,omitempty"`
	LinkInfo      *LinkInfo    `json:"linkInfo,omitempty"`
	Dev           *DevInfo     `json:"dev,omitempty"`
	TopHolderPct  float64      `json:"topHolderPct,omitempty"`
	TopHolderList []HolderInfo `json:"topHolderList,omitempty"`
}

// TradeInfo contains trading statistics for a token.
type TradeInfo struct {
	MarketCapUSD    float64 `json:"marketCapUsd"`
	Price           float64 `json:"price"`
	Holder          int     `json:"holder"`
	HourTradeNum    int     `json:"hourTradeNum"`
	HourTradeVolume float64 `json:"hourTradeVolume"`
}

// PairInfo contains trading pair information.
type PairInfo struct {
	PairAddress  string  `json:"pairAddress"`
	Pair         string  `json:"pair"`
	LiquidateUSD float64 `json:"liquidateUsd"`
	CreateTime   int64   `json:"createTime"`
}

// SecurityInfo contains security check results for a token.
type SecurityInfo struct {
	HoneyPot   bool `json:"honeyPot"`
	OpenSource bool `json:"openSource"`
	NoOwner    bool `json:"noOwner"`
	Locked     bool `json:"locked"`
}

// TaxInfo contains buy/sell tax rates (as percentage strings).
type TaxInfo struct {
	Buy  string `json:"buy"`
	Sell string `json:"sell"`
}

// LinkInfo contains social media links for a token.
type LinkInfo struct {
	TG  string `json:"tg"`
	X   string `json:"x"`
	Web string `json:"web"`
}

// DevInfo contains developer information.
type DevInfo struct {
	Address string  `json:"address"`
	Pct     float64 `json:"pct"`
}

// HolderInfo represents a top holder entry.
type HolderInfo struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
	Pct     float64 `json:"pct"`
}

// ---------------------------------------------------------------------------
// Wallets
// ---------------------------------------------------------------------------

// WalletsRequest represents the request parameters for listing wallets.
type WalletsRequest struct {
	Chain        Chain  `json:"chain"`
	PageNum      int    `json:"pageNum,omitempty"`
	PageSize     int    `json:"pageSize,omitempty"`
	TokenAddress string `json:"tokenAddress,omitempty"`
}

// WalletsData represents the response from listing wallets.
type WalletsData struct {
	TotalCount int          `json:"totalCount"`
	PageSize   int          `json:"pageSize"`
	TotalPage  int          `json:"totalPage"`
	CurrPage   int          `json:"currPage"`
	List       []WalletItem `json:"list"`
}

// WalletItem represents a single wallet in the list response.
type WalletItem struct {
	UserID       int           `json:"userId"`
	Chain        ChainCode     `json:"chain"`
	Name         string        `json:"name"`
	Address      string        `json:"address"`
	Balance      float64       `json:"balance"`
	TopUp        int           `json:"topUp"` // 1=pinned, 0=normal
	TokenBalance *TokenBalance `json:"tokenBalance,omitempty"`
	CreateTime   string        `json:"createTime"`
	UpdateTime   string        `json:"updateTime"`
	IsImport     bool          `json:"isImport"`
}

// IsPinned returns true if the wallet is pinned.
func (w *WalletItem) IsPinned() bool {
	return w.TopUp == 1
}

// WalletInfoRequest represents the request parameters for wallet info.
type WalletInfoRequest struct {
	WalletAddress string `json:"walletAddress"`
	Chain         Chain  `json:"chain"`
	TokenAddress  string `json:"tokenAddress,omitempty"`
}

// WalletInfoData represents the response from querying wallet info.
type WalletInfoData struct {
	Address      string        `json:"address"`
	Name         string        `json:"name"`
	Chain        ChainCode     `json:"chain"`
	IsImport     bool          `json:"isImport"`
	TopUp        int           `json:"topUp"`
	Balance      float64       `json:"balance"`
	TokenBalance *TokenBalance `json:"tokenBalance,omitempty"`
}

// IsPinned returns true if the wallet is pinned.
func (w *WalletInfoData) IsPinned() bool {
	return w.TopUp == 1
}

// TokenBalance represents token holdings in a wallet.
type TokenBalance struct {
	Amount         string  `json:"amount"`
	UIAmount       float64 `json:"uiAmount"`
	Decimals       int     `json:"decimals"`
	UIAmountString string  `json:"uiAmountString,omitempty"`
}
