package main

import (
	"fmt"

	"Dabber/binance"
)

func main() {
	b := binance.NewBinanceClient()

	r, err := b.GetTradeRatio("DENT", "BTC")
	if err != nil {
		panic(err)
	}

	fmt.Println(r)
}
