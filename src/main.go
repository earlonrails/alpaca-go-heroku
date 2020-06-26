package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
	"github.com/shopspring/decimal"
)

func init() {
	API_KEY := "YOUR_API_KEY_HERE"
	API_SECRET := "YOUR_API_SECRET_HERE"
	BASE_URL := "https://paper-api.alpaca.markets"

	// Check for environment variables
	if common.Credentials().ID == "" {
		os.Setenv(common.EnvApiKeyID, API_KEY)
	}
	if common.Credentials().Secret == "" {
		os.Setenv(common.EnvApiSecretKey, API_SECRET)
	}
	alpaca.SetBaseUrl(BASE_URL)

	// Format the allStocks variable for use in the class.
	allStocks := []stockField{}
	stockList := []string{"DOMO", "TLRY", "SQ", "MRO", "AAPL", "GM", "SNAP", "SHOP", "SPLK", "BA", "AMZN", "SUI", "SUN", "TSLA", "CGC", "SPWR", "NIO", "CAT", "MSFT", "PANW", "OKTA", "TWTR", "TM", "RTN", "ATVI", "GS", "BAC", "MS", "TWLO", "QCOM"}
	for _, stock := range stockList {
		allStocks = append(allStocks, stockField{stock, 0})
	}

	alpacaClient = alpacaClientContainer{
		alpaca.NewClient(common.Credentials()),
		bucket{[]string{}, -1, -1, 0},
		bucket{[]string{}, -1, -1, 0},
		make([]stockField, len(allStocks)),
		[]string{},
	}

	copy(alpacaClient.allStocks, allStocks)
}

func main() {
	// First, cancel any existing orders so they don't impact our buying power.
	status, until, limit := "open", time.Now(), 100
	orders, _ := alpacaClient.client.ListOrders(&status, &until, &limit, nil)
	for _, order := range orders {
		_ = alpacaClient.client.CancelOrder(order.ID)
	}

	// Wait for market to open.
	fmt.Println("Waiting for market to open...")
	for {
		isOpen := alpacaClient.awaitMarketOpen()
		if isOpen {
			break
		}
		time.Sleep(1 * time.Minute)
	}
	fmt.Println("Market Opened.")

	for {
		alpacaClient.run()
	}
}
