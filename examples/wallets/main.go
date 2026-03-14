package main

import (
	"context"
	"fmt"
	"os"

	xxyy "github.com/QuantProcessing/xxyy-api-go"
)

func main() {
	apiKey := os.Getenv("XXYY_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "XXYY_API_KEY environment variable is required")
		os.Exit(1)
	}

	client := xxyy.NewClient(apiKey)
	ctx := context.Background()

	// List wallets on Solana
	fmt.Println("=== Wallets on SOL ===")
	data, err := client.ListWallets(ctx, &xxyy.WalletsRequest{
		Chain: xxyy.ChainSOL,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ListWallets failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Total wallets: %d\n\n", data.TotalCount)
	for _, w := range data.List {
		pinned := ""
		if w.IsPinned() {
			pinned = "⭐ "
		}
		fmt.Printf("%s%s\n", pinned, w.Name)
		fmt.Printf("  Address: %s\n", w.Address)
		fmt.Printf("  Balance: %v %s\n\n", w.Balance, xxyy.NativeToken(xxyy.ChainSOL))
	}
}
