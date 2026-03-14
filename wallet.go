package xxyy

import (
	"context"
	"fmt"
	"strconv"
)

// ListWallets lists wallets on the XXYY platform for the current user.
//
// Returns wallet addresses and native token balances.
// Pass nil for req to use defaults (chain=sol, page 1, size 20).
func (c *Client) ListWallets(ctx context.Context, req *WalletsRequest) (*WalletsData, error) {
	params := make(map[string]string)

	if req != nil {
		if req.Chain != "" {
			if !req.Chain.IsValid() {
				return nil, fmt.Errorf("xxyy: invalid chain %q", req.Chain)
			}
			params["chain"] = string(req.Chain)
		}
		if req.PageNum > 0 {
			params["pageNum"] = strconv.Itoa(req.PageNum)
		}
		if req.PageSize > 0 {
			params["pageSize"] = strconv.Itoa(req.PageSize)
		}
		if req.TokenAddress != "" {
			params["tokenAddress"] = req.TokenAddress
		}
	}

	resp, err := doGet[WalletsData](ctx, c, apiPrefix+"/wallets", params)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, newXxyyError(resp.Code, fmt.Sprintf("wallet query failed: %s", resp.Msg))
	}

	return &resp.Data, nil
}

// WalletInfo queries a single wallet's details including native token balance
// and optional token holdings.
func (c *Client) WalletInfo(ctx context.Context, req *WalletInfoRequest) (*WalletInfoData, error) {
	if req == nil {
		return nil, fmt.Errorf("xxyy: wallet info request must not be nil")
	}
	if req.WalletAddress == "" {
		return nil, fmt.Errorf("xxyy: walletAddress must not be empty")
	}

	chain := req.Chain
	if chain == "" {
		chain = ChainSOL
	}
	if !chain.IsValid() {
		return nil, fmt.Errorf("xxyy: invalid chain %q", chain)
	}

	if err := ValidateWalletAddress(req.WalletAddress, chain); err != nil {
		return nil, fmt.Errorf("xxyy: %w", err)
	}

	params := map[string]string{
		"walletAddress": req.WalletAddress,
		"chain":         string(chain),
	}
	if req.TokenAddress != "" {
		params["tokenAddress"] = req.TokenAddress
	}

	resp, err := doGet[WalletInfoData](ctx, c, apiPrefix+"/wallet/info", params)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, newXxyyError(resp.Code, fmt.Sprintf("wallet info query failed: %s", resp.Msg))
	}

	return &resp.Data, nil
}
