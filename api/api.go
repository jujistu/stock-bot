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

type alphaVantageResponse struct {
	GlobalQuote struct {
		Symbol        string `json:"01. symbol"`
		Price         string `json:"05. price"`
		LatestTrading string `json:"07. latest trading day"`
	} `json:"Global Quote"`
}

func EvalStock(key string) string {
	symbol := strings.TrimSpace(strings.ToUpper(key))

	if symbol == "" {
		return "Stock symbol is required"
	}

	// Accept the existing command format, e.g. aapl.us,
	// while Alpha Vantage expects AAPL.
	symbol = strings.TrimSuffix(symbol, ".US")

	apiKey := os.Getenv("STOCK_API_KEY")
	if apiKey == "" {
		log.Println("error: STOCK_API_KEY is not configured")
		return "Stock service is not configured"
	}

	params := url.Values{}
	params.Set("function", "GLOBAL_QUOTE")
	params.Set("symbol", symbol)
	params.Set("apikey", apiKey)

	stockServiceURL := "https://www.alphavantage.co/query?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, stockServiceURL, nil)
	if err != nil {
		log.Println("error creating stock request:", err)
		return "Stock service request error"
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("error requesting stock service:", err)
		return "Stock service is not available"
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Println("error: response.StatusCode is", response.StatusCode)
		return "Stock service is not available"
	}

	var data alphaVantageResponse

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		log.Println("error decoding stock response:", err)
		return "Stock service response error"
	}

	if data.GlobalQuote.Symbol == "" || data.GlobalQuote.Price == "" {
		log.Println("error: stock quote unavailable for", symbol)
		return fmt.Sprintf("%s quote is not available", symbol)
	}

	return fmt.Sprintf(
		"%s quote is $%s per share",
		strings.ToUpper(data.GlobalQuote.Symbol),
		data.GlobalQuote.Price,
	)
}
