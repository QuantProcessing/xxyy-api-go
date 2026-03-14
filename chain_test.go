package xxyy

import (
	"testing"
)

func TestChainIsValid(t *testing.T) {
	tests := []struct {
		chain Chain
		want  bool
	}{
		{ChainSOL, true},
		{ChainETH, true},
		{ChainBSC, true},
		{ChainBase, true},
		{Chain("polygon"), false},
		{Chain(""), false},
	}
	for _, tt := range tests {
		if got := tt.chain.IsValid(); got != tt.want {
			t.Errorf("Chain(%q).IsValid() = %v, want %v", tt.chain, got, tt.want)
		}
	}
}

func TestChainIsFeedSupported(t *testing.T) {
	tests := []struct {
		chain Chain
		want  bool
	}{
		{ChainSOL, true},
		{ChainBSC, true},
		{ChainETH, false},
		{ChainBase, false},
	}
	for _, tt := range tests {
		if got := tt.chain.IsFeedSupported(); got != tt.want {
			t.Errorf("Chain(%q).IsFeedSupported() = %v, want %v", tt.chain, got, tt.want)
		}
	}
}

func TestExplorerURL(t *testing.T) {
	tests := []struct {
		chain Chain
		txID  string
		want  string
	}{
		{ChainSOL, "abc123", "https://solscan.io/tx/abc123"},
		{ChainETH, "0xdef456", "https://etherscan.io/tx/0xdef456"},
		{ChainBSC, "0xabc", "https://bscscan.com/tx/0xabc"},
		{ChainBase, "0xbase", "https://basescan.org/tx/0xbase"},
	}
	for _, tt := range tests {
		if got := ExplorerURL(tt.chain, tt.txID); got != tt.want {
			t.Errorf("ExplorerURL(%q, %q) = %q, want %q", tt.chain, tt.txID, got, tt.want)
		}
	}
}

func TestNativeToken(t *testing.T) {
	tests := []struct {
		chain Chain
		want  string
	}{
		{ChainSOL, "SOL"},
		{ChainETH, "ETH"},
		{ChainBSC, "BNB"},
		{ChainBase, "ETH"},
	}
	for _, tt := range tests {
		if got := NativeToken(tt.chain); got != tt.want {
			t.Errorf("NativeToken(%q) = %q, want %q", tt.chain, got, tt.want)
		}
	}
}

func TestChainCodeToChain(t *testing.T) {
	tests := []struct {
		code      ChainCode
		wantChain Chain
		wantOK    bool
	}{
		{ChainCodeSOL, ChainSOL, true},
		{ChainCodeBSC, ChainBSC, true},
		{ChainCodeETH, ChainETH, true},
		{ChainCodeBase, ChainBase, true},
		{ChainCode(99), "", false},
	}
	for _, tt := range tests {
		chain, ok := tt.code.ToChain()
		if ok != tt.wantOK || chain != tt.wantChain {
			t.Errorf("ChainCode(%d).ToChain() = (%q, %v), want (%q, %v)",
				tt.code, chain, ok, tt.wantChain, tt.wantOK)
		}
	}
}
