package xxyy

import (
	"testing"
)

func TestValidateWalletAddress(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		chain   Chain
		wantErr bool
	}{
		// Valid Solana addresses
		{"valid SOL 44 chars", "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU", ChainSOL, false},
		{"valid SOL 32 chars", "11111111111111111111111111111111", ChainSOL, false},

		// Invalid Solana addresses
		{"SOL too short", "abc", ChainSOL, true},
		{"SOL invalid chars (0)", "0xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAs", ChainSOL, true},
		{"SOL invalid chars (O)", "OxKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAs", ChainSOL, true},
		{"SOL invalid chars (I)", "IxKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAs", ChainSOL, true},
		{"SOL invalid chars (l)", "lxKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAs", ChainSOL, true},

		// Valid EVM addresses
		{"valid ETH", "0x1a2B3c4D5e6F7a8B9c0D1E2F3a4B5C6D7E8F9A0b", ChainETH, false},
		{"valid BSC", "0x1234567890abcdef1234567890abcdef12345678", ChainBSC, false},
		{"valid Base", "0xABCDEF1234567890ABCDEF1234567890ABCDEF12", ChainBase, false},

		// Invalid EVM addresses
		{"ETH missing 0x", "1234567890abcdef1234567890abcdef12345678", ChainETH, true},
		{"ETH too short", "0x1234", ChainETH, true},
		{"ETH too long", "0x1234567890abcdef1234567890abcdef1234567890", ChainETH, true},
		{"ETH invalid chars", "0x1234567890abcdef1234567890abcdef1234567g", ChainETH, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWalletAddress(tt.addr, tt.chain)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateWalletAddress(%q, %q) error = %v, wantErr %v",
					tt.addr, tt.chain, err, tt.wantErr)
			}
		})
	}
}

func TestValidateContractAddress(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		chain   Chain
		wantErr bool
	}{
		{"valid SOL contract", "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", ChainSOL, false},
		{"valid EVM contract", "0x1234567890abcdef1234567890abcdef12345678", ChainETH, false},
		{"invalid SOL", "invalid", ChainSOL, true},
		{"invalid EVM", "notanaddress", ChainBSC, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateContractAddress(tt.addr, tt.chain)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateContractAddress(%q, %q) error = %v, wantErr %v",
					tt.addr, tt.chain, err, tt.wantErr)
			}
		})
	}
}

func TestValidateTip(t *testing.T) {
	tests := []struct {
		name     string
		tip      float64
		chain    Chain
		wantWarn bool
	}{
		// SOL range: 0.001-0.1
		{"SOL valid min", 0.001, ChainSOL, false},
		{"SOL valid mid", 0.05, ChainSOL, false},
		{"SOL valid max", 0.1, ChainSOL, false},
		{"SOL too low", 0.0001, ChainSOL, true},
		{"SOL too high", 0.5, ChainSOL, true},

		// EVM range: 0.1-100
		{"ETH valid min", 0.1, ChainETH, false},
		{"ETH valid mid", 50, ChainETH, false},
		{"ETH valid max", 100, ChainETH, false},
		{"ETH too low", 0.01, ChainETH, true},
		{"ETH too high", 200, ChainETH, true},

		{"BSC valid", 1, ChainBSC, false},
		{"Base valid", 5, ChainBase, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warning := ValidateTip(tt.tip, tt.chain)
			hasWarn := warning != ""
			if hasWarn != tt.wantWarn {
				t.Errorf("ValidateTip(%v, %q) warning = %q, wantWarn %v",
					tt.tip, tt.chain, warning, tt.wantWarn)
			}
		})
	}
}

func TestValidateSwapParams(t *testing.T) {
	slippage20 := 20
	model1 := 1
	model3 := 3
	priorityFee := 0.001

	validSolReq := &SwapRequest{
		Chain:         ChainSOL,
		WalletAddress: "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU",
		TokenAddress:  "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		Amount:        0.1,
		Tip:           0.001,
		Slippage:      &slippage20,
		Model:         &model1,
	}

	t.Run("valid buy", func(t *testing.T) {
		if err := validateSwapParams(validSolReq, true); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("valid sell", func(t *testing.T) {
		req := *validSolReq
		req.Amount = 50
		if err := validateSwapParams(&req, false); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("invalid chain", func(t *testing.T) {
		req := *validSolReq
		req.Chain = "polygon"
		if err := validateSwapParams(&req, true); err == nil {
			t.Error("expected error for invalid chain")
		}
	})

	t.Run("buy amount <= 0", func(t *testing.T) {
		req := *validSolReq
		req.Amount = 0
		if err := validateSwapParams(&req, true); err == nil {
			t.Error("expected error for amount <= 0")
		}
	})

	t.Run("sell amount out of range", func(t *testing.T) {
		req := *validSolReq
		req.Amount = 101
		if err := validateSwapParams(&req, false); err == nil {
			t.Error("expected error for sell amount > 100")
		}
	})

	t.Run("invalid model", func(t *testing.T) {
		req := *validSolReq
		req.Model = &model3
		if err := validateSwapParams(&req, true); err == nil {
			t.Error("expected error for model = 3")
		}
	})

	t.Run("priorityFee on EVM", func(t *testing.T) {
		req := &SwapRequest{
			Chain:         ChainETH,
			WalletAddress: "0x1234567890abcdef1234567890abcdef12345678",
			TokenAddress:  "0xabcdef1234567890abcdef1234567890abcdef12",
			Amount:        0.1,
			Tip:           1,
			PriorityFee:   &priorityFee,
		}
		if err := validateSwapParams(req, true); err == nil {
			t.Error("expected error for priorityFee on EVM chain")
		}
	})

	t.Run("priorityFee on SOL allowed", func(t *testing.T) {
		req := *validSolReq
		req.PriorityFee = &priorityFee
		if err := validateSwapParams(&req, true); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
