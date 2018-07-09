package main

import (
	"encoding/json"
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
		b, err := json.MarshalIndent(ts, "", "")
		if err != nil {
			panic(err)
		}

		fmt.Println(string(b))

	}

}
