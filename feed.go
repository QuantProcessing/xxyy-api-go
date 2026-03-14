package xxyy

import (
	"context"
	"fmt"
)

// FeedScan retrieves Meme token lists: newly launched, almost graduated, or graduated.
//
// Only SOL and BSC chains are supported for feed scanning.
// Pass nil for filter to use no filters (return all tokens).
func (c *Client) FeedScan(ctx context.Context, feedType FeedType, chain Chain, filter *FeedFilter) (*FeedData, error) {
	if !feedType.IsValid() {
		return nil, fmt.Errorf("xxyy: invalid feed type %q: must be NEW, ALMOST, or COMPLETED", feedType)
	}
	if !chain.IsFeedSupported() {
		return nil, fmt.Errorf("xxyy: feed scanning only supports sol and bsc chains, got %q", chain)
	}

	// Use empty object if no filter provided (same as TypeScript implementation)
	var body any
	if filter != nil {
		body = filter
	} else {
		body = struct{}{}
	}

	path := fmt.Sprintf("%s/feed/%s", apiPrefix, feedType)
	params := map[string]string{"chain": string(chain)}

	resp, err := doPost[FeedData](ctx, c, path, body, params)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, newXxyyError(resp.Code, fmt.Sprintf("feed scan failed: %s", resp.Msg))
	}

	return &resp.Data, nil
}
