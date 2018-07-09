package binance

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	webSocketURL = "wss://stream.binance.com:9443"
)

type tickerStreamInfo struct {
	eventTime       time.Time
	bestBidPrice    float64
	bestBidQuantity float64
	bestAskPrice    float64
	bestAskQuantity float64
}

// GetTickerStream returns a channel on which the caller can receive ticker info.
func (bc *binanceClient) GetTickerStream(symbol string) (<-chan tickerStreamInfo, error) {
	ticker := make(chan tickerStreamInfo)

	addr := fmt.Sprintf("%s/ws/%s@ticker", webSocketURL, symbol)

	var wsDialer websocket.Dialer
	conn, _, err := wsDialer.Dial(addr, nil)
	if err != nil {
		return ticker, err
	}

	bc.wg.Add(1)
	go func() {
		defer bc.wg.Done()
		defer close(ticker)
		defer conn.Close()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("could not read from stream: %v\n", err)
				return
			}

			ts, err := extractTickerStreamInfo(msg)
			if err != nil {
				fmt.Printf("could not extract ticker stream info: %v\n", err)
				return
			}
			ticker <- ts
			time.Sleep(500 * time.Millisecond)
		}
	}()

	return ticker, nil
}

func (ts tickerStreamInfo) String() string {
	return fmt.Sprintf("%v: best bid price: %v at %v,  best ask price: %v at %v",
		ts.eventTime, ts.bestBidPrice, ts.bestBidQuantity, ts.bestAskPrice, ts.bestAskQuantity)
}

// extractTickerStreamInfo converts a raw websocket message into a tickerStreamInfo.
// see https://github.com/binance-exchange/binance-official-api-docs/blob/master/web-socket-streams.md#individual-symbol-ticker-streams
func extractTickerStreamInfo(raw []byte) (tickerStreamInfo, error) {
	var ts tickerStreamInfo
	var jsonResp struct {
		//ET  int    `json:"E"`
		BBP string `json:"b"`
		BBQ string `json:"B"`
		BAP string `json:"a"`
		BAQ string `json:"A"`
	}

	if err := json.Unmarshal(raw, &jsonResp); err != nil {
		return ts, err
	}

	// fmt.Println(string(raw))
	//ts.eventTime = time.Unix(0, int64(jsonResp.ET*int(time.Millisecond)))

	bbp, err := strconv.ParseFloat(jsonResp.BBP, 64)
	if err != nil {
		return ts, err
	}
	ts.bestBidPrice = bbp

	bbq, err := strconv.ParseFloat(jsonResp.BBQ, 64)
	if err != nil {
		return ts, err
	}
	ts.bestBidQuantity = bbq

	bap, err := strconv.ParseFloat(jsonResp.BAP, 64)
	if err != nil {
		return ts, err
	}
	ts.bestAskPrice = bap

	baq, err := strconv.ParseFloat(jsonResp.BAQ, 64)
	if err != nil {
		return ts, err
	}
	ts.bestAskQuantity = baq

	return ts, nil
}
