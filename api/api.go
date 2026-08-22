package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type AlphaVantageResponse struct {
	GlobalQuote struct {
		Symbol        string `json:"01. symbol"`
		Open          string `json:"02. open"`
		High          string `json:"03. high"`
		Low           string `json:"04. low"`
		Price         string `json:"05. price"`
		Volume        string `json:"06. volume"`
		LatestTrading string `json:"07. latest trading day"`
		PreviousClose string `json:"08. previous close"`
		Change        string `json:"09. change"`
		ChangePercent string `json:"10. change percent"`
	} `json:"Global Quote"`
}

type AlphaVantageErrorResponse struct {
	Information string `json:"Information"`
	Note        string `json:"Note"`
	Error       string `json:"Error Message"`
}

func EvalStock(key string) string {
	symbol := strings.TrimSpace(key)

	if symbol == "" {
		return "Stock symbol is required"
	}

	apiKey := os.Getenv("STOCK_API_KEY")
	if apiKey == "" {
		log.Println("error: STOCK_API_KEY is not configured")
		return "Stock service is not configured"
	}

	// The old application accepted symbols such as aapl.us.
	// Alpha Vantage uses AAPL for US equities.
	symbol = strings.TrimSuffix(strings.ToUpper(symbol), ".US")

	stockServiceURL := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s",
		url.QueryEscape(symbol),
		url.QueryEscape(apiKey),
	)

	log.Println("info: processing", stockServiceURL)

	response, err := http.Get(stockServiceURL)
	if err != nil {
		log.Println("error:", err)
		return "Stock service is not available"
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Println("error: response.StatusCode is", response.StatusCode)
		return "Stock service is not available"
	}

	var result AlphaVantageResponse

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Println("error decoding stock response:", err)
		return "Stock service response error"
	}

	if result.GlobalQuote.Symbol == "" {
		log.Println("error: no quote returned for", symbol)
		return fmt.Sprintf("%s quote is not available", symbol)
	}

	price := result.GlobalQuote.Price

	if price == "" {
		return fmt.Sprintf("%s quote is not available", symbol)
	}

	return fmt.Sprintf(
		"%s quote is $%s per share",
		strings.ToUpper(result.GlobalQuote.Symbol),
		price,
	)
}
