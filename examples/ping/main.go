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
		fmt.Fprintln(os.Stderr, "Get your API key at https://www.xxyy.io/apikey")
		os.Exit(1)
	}

	client := xxyy.NewClient(apiKey)

	if err := client.Ping(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Ping failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("pong — API Key is valid.")
}
