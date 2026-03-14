package xxyy

// Chain represents a supported blockchain network.
type Chain string

const (
	ChainSOL  Chain = "sol"
	ChainETH  Chain = "eth"
	ChainBSC  Chain = "bsc"
	ChainBase Chain = "base"
)

// ValidChains lists all supported chains.
var ValidChains = []Chain{ChainSOL, ChainETH, ChainBSC, ChainBase}

// FeedChains lists chains that support the Feed API.
var FeedChains = []Chain{ChainSOL, ChainBSC}

// IsValid returns true if the chain is a supported chain.
func (c Chain) IsValid() bool {
	switch c {
	case ChainSOL, ChainETH, ChainBSC, ChainBase:
		return true
	}
	return false
}

// IsFeedSupported returns true if the chain supports the Feed API.
func (c Chain) IsFeedSupported() bool {
	return c == ChainSOL || c == ChainBSC
}

// String returns the string representation of the chain.
func (c Chain) String() string {
	return string(c)
}

var explorerURLs = map[Chain]string{
	ChainSOL:  "https://solscan.io/tx/",
	ChainETH:  "https://etherscan.io/tx/",
	ChainBSC:  "https://bscscan.com/tx/",
	ChainBase: "https://basescan.org/tx/",
}

// ExplorerURL returns the block explorer URL for a transaction.
func ExplorerURL(chain Chain, txID string) string {
	prefix, ok := explorerURLs[chain]
	if !ok {
		return txID
	}
	return prefix + txID
}

var nativeTokens = map[Chain]string{
	ChainSOL:  "SOL",
	ChainETH:  "ETH",
	ChainBSC:  "BNB",
	ChainBase: "ETH",
}

// NativeToken returns the native token symbol for a chain.
func NativeToken(chain Chain) string {
	if t, ok := nativeTokens[chain]; ok {
		return t
	}
	return string(chain)
}

// ChainCode represents the numeric chain identifier used in API responses.
type ChainCode int

const (
	ChainCodeSOL  ChainCode = 1
	ChainCodeBSC  ChainCode = 2
	ChainCodeETH  ChainCode = 3
	ChainCodeBase ChainCode = 6
)

var chainCodeMap = map[ChainCode]Chain{
	ChainCodeSOL:  ChainSOL,
	ChainCodeBSC:  ChainBSC,
	ChainCodeETH:  ChainETH,
	ChainCodeBase: ChainBase,
}

// ToChain converts a numeric chain code to a Chain value.
func (cc ChainCode) ToChain() (Chain, bool) {
	c, ok := chainCodeMap[cc]
	return c, ok
}
