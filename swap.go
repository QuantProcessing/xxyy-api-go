package xxyy

import (
	"context"
	"fmt"
)

// BuyToken submits a buy order for a token.
//
// Amount is in native currency (SOL/ETH/BNB).
// The request is NOT retried on rate limit since it's an irreversible operation.
func (c *Client) BuyToken(ctx context.Context, req *SwapRequest) (*SwapResponse, error) {
	return c.doSwap(ctx, true, req)
}

// SellToken submits a sell order for a token.
//
// Amount is a sell percentage (1-100). Example: 50 = sell 50% of holdings.
// The request is NOT retried on rate limit since it's an irreversible operation.
func (c *Client) SellToken(ctx context.Context, req *SwapRequest) (*SwapResponse, error) {
	return c.doSwap(ctx, false, req)
}

func (c *Client) doSwap(ctx context.Context, isBuy bool, req *SwapRequest) (*SwapResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("xxyy: swap request must not be nil")
	}

	// Validate parameters
	if err := validateSwapParams(req, isBuy); err != nil {
		return nil, fmt.Errorf("xxyy: %w", err)
	}

	// Build request body — only include priorityFee for SOL chain
	body := swapBody{
		Chain:         req.Chain,
		WalletAddress: req.WalletAddress,
		TokenAddress:  req.TokenAddress,
		IsBuy:         isBuy,
		Amount:        req.Amount,
		Tip:           req.Tip,
		Slippage:      req.Slippage,
		Model:         req.Model,
	}
	if req.PriorityFee != nil && req.Chain == ChainSOL {
		body.PriorityFee = req.PriorityFee
	}

	resp, err := doPostNoRetry[SwapResponse](ctx, c, apiPrefix+"/swap", body, nil)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, newXxyyError(resp.Code, fmt.Sprintf("swap failed: %s", resp.Msg))
	}

	return &resp.Data, nil
}
