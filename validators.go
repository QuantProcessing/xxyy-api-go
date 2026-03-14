package xxyy

import (
	"fmt"
	"regexp"
)

var (
	solAddressRE = regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]{32,44}$`)
	evmAddressRE = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
)

type tipRange struct {
	Min  float64
	Max  float64
	Unit string
}

var tipRanges = map[Chain]tipRange{
	ChainSOL:  {Min: 0.001, Max: 0.1, Unit: "SOL"},
	ChainETH:  {Min: 0.1, Max: 100, Unit: "Gwei"},
	ChainBSC:  {Min: 0.1, Max: 100, Unit: "Gwei"},
	ChainBase: {Min: 0.1, Max: 100, Unit: "Gwei"},
}

// ValidateWalletAddress validates a wallet address format for the given chain.
func ValidateWalletAddress(address string, chain Chain) error {
	if chain == ChainSOL {
		if !solAddressRE.MatchString(address) {
			return fmt.Errorf("invalid Solana wallet address: expected Base58, 32-44 characters, got %q", address)
		}
		return nil
	}
	if !evmAddressRE.MatchString(address) {
		return fmt.Errorf("invalid EVM wallet address for %s: expected 0x + 40 hex characters, got %q", chain, address)
	}
	return nil
}

// ValidateContractAddress validates a contract address format for the given chain.
func ValidateContractAddress(address string, chain Chain) error {
	if chain == ChainSOL {
		if !solAddressRE.MatchString(address) {
			return fmt.Errorf("invalid Solana contract address: expected Base58, 32-44 characters, got %q", address)
		}
		return nil
	}
	if !evmAddressRE.MatchString(address) {
		return fmt.Errorf("invalid EVM contract address for %s: expected 0x + 40 hex characters, got %q", chain, address)
	}
	return nil
}

// ValidateTip checks if the tip value is within the recommended range for the chain.
// Returns a warning string if the tip is outside the recommended range, or empty string if OK.
func ValidateTip(tip float64, chain Chain) string {
	r, ok := tipRanges[chain]
	if !ok {
		return ""
	}
	if tip < r.Min || tip > r.Max {
		return fmt.Sprintf(
			"tip %.6g is outside the recommended range (%.6g-%.6g %s for %s); this may result in unexpectedly high costs or failed transactions",
			tip, r.Min, r.Max, r.Unit, chain,
		)
	}
	return ""
}

// validateSwapParams performs full validation on swap parameters, matching the
// TypeScript implementation's strict parameter validation (Execution Rules #8).
func validateSwapParams(req *SwapRequest, isBuy bool) error {
	if !req.Chain.IsValid() {
		return fmt.Errorf("invalid chain %q: must be one of sol/eth/bsc/base", req.Chain)
	}
	if err := ValidateWalletAddress(req.WalletAddress, req.Chain); err != nil {
		return err
	}
	if err := ValidateContractAddress(req.TokenAddress, req.Chain); err != nil {
		return err
	}
	if isBuy {
		if req.Amount <= 0 {
			return fmt.Errorf("buy amount must be > 0, got %v", req.Amount)
		}
	} else {
		if req.Amount < 1 || req.Amount > 100 {
			return fmt.Errorf("sell amount (percentage) must be 1-100, got %v", req.Amount)
		}
	}
	if warning := ValidateTip(req.Tip, req.Chain); warning != "" {
		return fmt.Errorf("%s", warning)
	}
	if req.Slippage != nil {
		if *req.Slippage < 0 || *req.Slippage > 100 {
			return fmt.Errorf("slippage must be 0-100, got %v", *req.Slippage)
		}
	}
	if req.Model != nil {
		if *req.Model != 1 && *req.Model != 2 {
			return fmt.Errorf("model must be 1 or 2, got %v", *req.Model)
		}
	}
	if req.PriorityFee != nil && req.Chain != ChainSOL {
		return fmt.Errorf("priorityFee is only supported on Solana chain")
	}
	return nil
}
