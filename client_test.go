package xxyy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// newTestServer creates a test HTTP server that returns the given response.
func newTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	client := NewClient("test_api_key", WithBaseURL(ts.URL))
	return ts, client
}

func jsonResponse(t *testing.T, w http.ResponseWriter, resp any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		t.Fatalf("failed to encode response: %v", err)
	}
}

func TestPing(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/trade/open/api/ping" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test_api_key" {
			t.Errorf("unexpected auth header: %s", r.Header.Get("Authorization"))
		}
		jsonResponse(t, w, map[string]any{
			"code":    200,
			"msg":     "success",
			"data":    "pong",
			"success": true,
		})
	})

	err := client.Ping(context.Background())
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
}

func TestQueryTrade(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("txId") != "test_tx_123" {
			t.Errorf("unexpected txId: %s", r.URL.Query().Get("txId"))
		}
		jsonResponse(t, w, map[string]any{
			"code": 200,
			"msg":  "success",
			"data": map[string]any{
				"txId":          "test_tx_123",
				"status":        "success",
				"statusDesc":    "Transaction successful",
				"chain":         "sol",
				"tokenAddress":  "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
				"walletAddress": "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU",
				"isBuy":         true,
				"baseAmount":    100.5,
				"quoteAmount":   0.1,
			},
			"success": true,
		})
	})

	data, err := client.QueryTrade(context.Background(), "test_tx_123")
	if err != nil {
		t.Fatalf("QueryTrade failed: %v", err)
	}

	if data.TxID != "test_tx_123" {
		t.Errorf("TxID = %q, want %q", data.TxID, "test_tx_123")
	}
	if data.Status != TradeStatusSuccess {
		t.Errorf("Status = %d, want %d (success)", data.Status, TradeStatusSuccess)
	}
	if !data.IsBuy {
		t.Error("IsBuy should be true")
	}
}

func TestAPIKeyError(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(t, w, map[string]any{
			"code":    8060,
			"msg":     "API Key invalid",
			"data":    nil,
			"success": false,
		})
	})

	err := client.Ping(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid API key")
	}

	xxErr, ok := err.(*XxyyError)
	if !ok {
		t.Fatalf("expected *XxyyError, got %T", err)
	}
	if !xxErr.IsAPIKeyError() {
		t.Errorf("expected IsAPIKeyError() = true, code = %d", xxErr.Code)
	}
}

func TestRateLimitRetry(t *testing.T) {
	var callCount atomic.Int32

	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		n := callCount.Add(1)
		if n <= 2 {
			jsonResponse(t, w, map[string]any{
				"code":    8062,
				"msg":     "rate limited",
				"data":    nil,
				"success": false,
			})
			return
		}
		jsonResponse(t, w, map[string]any{
			"code":    200,
			"msg":     "success",
			"data":    "pong",
			"success": true,
		})
	})

	// Override timeout for faster test
	client.timeout = 30 * time.Second

	err := client.Ping(context.Background())
	if err != nil {
		t.Fatalf("Ping should succeed after retries: %v", err)
	}

	if count := callCount.Load(); count != 3 {
		t.Errorf("expected 3 calls (2 retries + 1 success), got %d", count)
	}
}

func TestRateLimitExhausted(t *testing.T) {
	var callCount atomic.Int32

	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		callCount.Add(1)
		jsonResponse(t, w, map[string]any{
			"code":    8062,
			"msg":     "rate limited",
			"data":    nil,
			"success": false,
		})
	})

	err := client.Ping(context.Background())
	if err == nil {
		t.Fatal("expected error when retries exhausted")
	}

	xxErr, ok := err.(*XxyyError)
	if !ok {
		t.Fatalf("expected *XxyyError, got %T", err)
	}
	if !xxErr.IsRateLimited() {
		t.Errorf("expected IsRateLimited() = true, code = %d", xxErr.Code)
	}

	// Should have made 3 calls total (1 initial + 2 retries)
	if count := callCount.Load(); count != 3 {
		t.Errorf("expected 3 calls, got %d", count)
	}
}

func TestServerError(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(t, w, map[string]any{
			"code":    300,
			"msg":     "internal server error",
			"data":    nil,
			"success": false,
		})
	})

	err := client.Ping(context.Background())
	if err == nil {
		t.Fatal("expected error for server error")
	}

	xxErr, ok := err.(*XxyyError)
	if !ok {
		t.Fatalf("expected *XxyyError, got %T", err)
	}
	if !xxErr.IsServerError() {
		t.Errorf("expected IsServerError() = true, code = %d", xxErr.Code)
	}
}

func TestHTTPError(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	err := client.Ping(context.Background())
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}

	xxErr, ok := err.(*XxyyError)
	if !ok {
		t.Fatalf("expected *XxyyError, got %T", err)
	}
	if xxErr.Code != 500 {
		t.Errorf("expected code 500, got %d", xxErr.Code)
	}
}

func TestSwapNoRetry(t *testing.T) {
	var callCount atomic.Int32

	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		callCount.Add(1)
		jsonResponse(t, w, map[string]any{
			"code":    8062,
			"msg":     "rate limited",
			"data":    nil,
			"success": false,
		})
	})

	_, err := client.BuyToken(context.Background(), &SwapRequest{
		Chain:         ChainSOL,
		WalletAddress: "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU",
		TokenAddress:  "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		Amount:        0.1,
		Tip:           0.001,
	})

	if err == nil {
		t.Fatal("expected error for rate limited swap")
	}

	// Swap should NOT retry — only 1 call
	if count := callCount.Load(); count != 1 {
		t.Errorf("swap should not retry: expected 1 call, got %d", count)
	}
}

func TestListWallets(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/trade/open/api/wallets" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("chain") != "sol" {
			t.Errorf("unexpected chain: %s", r.URL.Query().Get("chain"))
		}
		jsonResponse(t, w, map[string]any{
			"code": 200,
			"msg":  "success",
			"data": map[string]any{
				"totalCount": 2,
				"pageSize":   20,
				"totalPage":  1,
				"currPage":   1,
				"list": []map[string]any{
					{
						"userId":     12345,
						"chain":      1,
						"name":       "Wallet-1",
						"address":    "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU",
						"balance":    1.5,
						"topUp":      1,
						"isImport":   false,
						"createTime": "2025-01-01 00:00:00",
						"updateTime": "2025-06-01 12:00:00",
					},
					{
						"userId":     12345,
						"chain":      1,
						"name":       "Wallet-2",
						"address":    "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
						"balance":    0.5,
						"topUp":      0,
						"isImport":   true,
						"createTime": "2025-02-01 00:00:00",
						"updateTime": "2025-06-01 12:00:00",
					},
				},
			},
			"success": true,
		})
	})

	data, err := client.ListWallets(context.Background(), &WalletsRequest{Chain: ChainSOL})
	if err != nil {
		t.Fatalf("ListWallets failed: %v", err)
	}

	if data.TotalCount != 2 {
		t.Errorf("TotalCount = %d, want 2", data.TotalCount)
	}
	if len(data.List) != 2 {
		t.Fatalf("List length = %d, want 2", len(data.List))
	}
	if !data.List[0].IsPinned() {
		t.Error("First wallet should be pinned")
	}
	if data.List[1].IsPinned() {
		t.Error("Second wallet should not be pinned")
	}
	if !data.List[1].IsImport {
		t.Error("Second wallet should be imported")
	}
}

func TestFeedScan(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/trade/open/api/feed/NEW" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		jsonResponse(t, w, map[string]any{
			"code": 200,
			"msg":  "success",
			"data": map[string]any{
				"items": []map[string]any{
					{
						"tokenAddress":   "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
						"symbol":         "TEST",
						"name":           "Test Token",
						"createTime":     1773140232851,
						"holders":        100,
						"priceUSD":       0.001,
						"marketCapUSD":   50000,
						"devHoldPercent": 5.5,
					},
				},
			},
			"success": true,
		})
	})

	data, err := client.FeedScan(context.Background(), FeedNew, ChainSOL, nil)
	if err != nil {
		t.Fatalf("FeedScan failed: %v", err)
	}

	if len(data.Items) != 1 {
		t.Fatalf("Items length = %d, want 1", len(data.Items))
	}
	if data.Items[0].Symbol != "TEST" {
		t.Errorf("Symbol = %q, want %q", data.Items[0].Symbol, "TEST")
	}
}

func TestFeedScanInvalidChain(t *testing.T) {
	client := NewClient("test_key")
	_, err := client.FeedScan(context.Background(), FeedNew, ChainETH, nil)
	if err == nil {
		t.Error("expected error for unsupported feed chain")
	}
}

func TestContextCancellation(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		jsonResponse(t, w, map[string]any{
			"code": 200, "msg": "success", "data": "pong", "success": true,
		})
	})

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := client.Ping(ctx)
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
}
