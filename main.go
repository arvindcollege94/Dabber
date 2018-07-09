package main

import (
	"fmt"

	"Dabber/binance"
)

func main() {
	b := binance.NewBinanceClient()
	defer b.Close()

	c, err := b.GetTickerStream("ethbtc")
	if err != nil {
		panic(err)
	}

	for ts := range c {
		fmt.Println(ts)
	}

}
