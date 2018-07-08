package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	binanceURL           = "https://api.binance.com"
	symbolPriceTickerURL = "api/v3/ticker/price"
)

type binanceClient struct {
	client *http.Client
}

func NewBinanceClient() *binanceClient {
	return &binanceClient{
		client: &http.Client{},
	}
}

// GetTradeRatio returns the ratio of b the exchange will give for a.
func (bc *binanceClient) GetTradeRatio(a, b string) (float64, error) {
	url := fmt.Sprintf("%s/%s?symbol=%s%s", binanceURL, symbolPriceTickerURL, a, b)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := bc.client.Do(req)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var jsonResp struct {
		Symbol string `json:"symbol"`
		// binance returns the price quoted so encoding/json
		// treats it as a string.
		Price string `json:"price"`
	}

	if err := json.Unmarshal(raw, &jsonResp); err != nil {
		return 0, err
	}

	return strconv.ParseFloat(jsonResp.Price, 64)
}
