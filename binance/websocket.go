package binance

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	webSocketURL = "wss://stream.binance.com:9443"
)

type tickerStreamInfo struct {
	eventTime       time.Time
	bestBidPrice    float64
	bestBidQuantity int
	bestAskPrice    float64
	bestAskQuantity int
}

// GetTickerStream returns a channel on which the caller can receive ticker info.
func (bc *binanceClient) GetTickerStream(symbol string) (<-chan tickerStreamInfo, error) {
	ret := make(chan tickerStreamInfo)

	addr := fmt.Sprintf("%s/ws/%s@ticker", webSocketURL, symbol)

	var wsDialer websocket.Dialer
	conn, _, err := wsDialer.Dial(addr, nil)
	if err != nil {
		return ret, err
	}
	defer conn.Close()

	bc.wg.Add(1)
	go func() {
		defer bc.wg.Done()
		defer close(ret)
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("could not read from stream: %v", err)
				return
			}

			var ts tickerStreamInfo
			if err := json.Unmarshal(msg, &ts); err != nil {
				fmt.Printf("could not unmarshal ticker stream message: %v", err)
				return
			}
			ret <- ts
			time.Sleep(500 * time.Millisecond)
		}
	}()

	return ret, nil
}
