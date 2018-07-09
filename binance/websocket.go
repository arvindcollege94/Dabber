package binance

import (
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

func (bc *binanceClient) GetTickerStream(symbol string) (<-chan tickerStreamInfo, error) {
	ret := make(chan tickerStreamInfo)

	addr := fmt.Sprintf("%s/ws/%s@ticker", webSocketURL, symbol)

	var wsDialer websocket.Dialer
	_, _, err := wsDialer.Dial(addr, nil)
	if err != nil {
		return ret, err
	}
	return ret, nil
}
