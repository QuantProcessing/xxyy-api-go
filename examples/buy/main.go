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

	// Replace these with actual values
	walletAddress := os.Getenv("XXYY_WALLET_ADDRESS")
	tokenAddress := os.Getenv("XXYY_TOKEN_ADDRESS")
	if walletAddress == "" || tokenAddress == "" {
		fmt.Fprintln(os.Stderr, "XXYY_WALLET_ADDRESS and XXYY_TOKEN_ADDRESS are required")
		os.Exit(1)
	}

	client := xxyy.NewClient(apiKey)
	ctx := context.Background()

	// Buy 0.001 SOL of a token
	slippage := 20
	resp, err := client.BuyToken(ctx, &xxyy.SwapRequest{
		Chain:         xxyy.ChainSOL,
		WalletAddress: walletAddress,
		TokenAddress:  tokenAddress,
		Amount:        0.001,
		Tip:           0.001,
		Slippage:      &slippage,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "BuyToken failed: %v\n", err)
		os.Exit(1)
	}

	txID := resp.TxID
	if txID == "" {
		txID = resp.Signature
	}

	fmt.Println("Buy order submitted!")
	fmt.Printf("Transaction ID: %s\n", txID)
	fmt.Printf("Explorer: %s\n", xxyy.ExplorerURL(xxyy.ChainSOL, txID))

	// Query trade status
	fmt.Println("\nQuerying trade status...")
	trade, err := client.QueryTrade(ctx, txID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryTrade failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Status: %s (%s)\n", trade.StatusText(), trade.StatusDesc)
}
