package xxyy

import (
	"context"
	"fmt"
)

// TokenQuery queries detailed information for a token: price, security checks,
// tax rates, holder distribution, and social links.
//
// Supports all 4 chains (sol/eth/bsc/base). Default chain is sol.
func (c *Client) TokenQuery(ctx context.Context, contractAddr string, chain Chain) (*TokenData, error) {
	if contractAddr == "" {
		return nil, fmt.Errorf("xxyy: contract address (ca) must not be empty")
	}
	if !chain.IsValid() {
		return nil, fmt.Errorf("xxyy: invalid chain %q: must be one of sol/eth/bsc/base", chain)
	}
	if err := ValidateContractAddress(contractAddr, chain); err != nil {
		return nil, fmt.Errorf("xxyy: %w", err)
	}

	resp, err := doGet[TokenData](ctx, c, apiPrefix+"/query", map[string]string{
		"ca":    contractAddr,
		"chain": string(chain),
	})
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, newXxyyError(resp.Code, fmt.Sprintf("token query failed: %s", resp.Msg))
	}

	return &resp.Data, nil
}
