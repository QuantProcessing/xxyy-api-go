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

	// Scan newly launched tokens on SOL with market cap > 10000 and holder > 50
	fmt.Println("=== Feed Scan: NEW tokens on SOL ===")
	data, err := client.FeedScan(ctx, xxyy.FeedNew, xxyy.ChainSOL, &xxyy.FeedFilter{
		MC:     "10000,",
		Holder: "50,",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "FeedScan failed: %v\n", err)
		os.Exit(1)
	}

	if len(data.Items) == 0 {
		fmt.Println("No tokens found matching filters.")
		return
	}

	fmt.Printf("Found %d tokens:\n\n", len(data.Items))
	for i, item := range data.Items {
		if i >= 10 { // Show first 10 only
			fmt.Printf("... and %d more\n", len(data.Items)-10)
			break
		}
		fmt.Printf("%s (%s)\n", item.Symbol, item.Name)
		fmt.Printf("  Contract: %s\n", item.TokenAddress)
		fmt.Printf("  Price:    $%v\n", item.PriceUSD)
		fmt.Printf("  MCap:     $%.0f\n", item.MarketCapUSD)
		fmt.Printf("  Holders:  %d\n", item.Holders)
		if item.LaunchPlatform != nil {
			fmt.Printf("  Platform: %s (progress: %s%%)\n",
				item.LaunchPlatform.Name, item.LaunchPlatform.Progress)
		}
		fmt.Println()
	}
}
