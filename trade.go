package xxyy

import (
	"context"
	"fmt"
)

// QueryTrade queries the status of a transaction by its ID.
func (c *Client) QueryTrade(ctx context.Context, txID string) (*TradeData, error) {
	if txID == "" {
		return nil, fmt.Errorf("xxyy: txId must not be empty")
	}

	resp, err := doGet[TradeData](ctx, c, apiPrefix+"/trade", map[string]string{
		"txId": txID,
	})
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, newXxyyError(resp.Code, fmt.Sprintf("query trade failed: %s", resp.Msg))
	}

	return &resp.Data, nil
}

// Ping verifies the XXYY API Key validity.
// Returns nil if the key is valid.
func (c *Client) Ping(ctx context.Context) error {
	resp, err := doGet[any](ctx, c, apiPrefix+"/ping", nil)
	if err != nil {
		return err
	}

	if !resp.Success {
		return newXxyyError(resp.Code, fmt.Sprintf("ping failed: %s", resp.Msg))
	}

	return nil
}
