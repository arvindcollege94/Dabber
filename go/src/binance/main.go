package main

import (
	"fmt"
)

func main() {
	b := NewBinanceClient()

	r, err := b.GetTradeRatio("DENT", "BTC")
	if err != nil {
		panic(err)
	}

	fmt.Println(r)
}
