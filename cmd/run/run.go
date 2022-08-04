package main

import (
	"fmt"

	prices "github.com/aanzolaavila/lapieza.io/internal"
)

func main() {
	file, err := prices.GetInputFile()
	if err != nil {
		panic(err)
	}

	outCh, errs := prices.ProcessPricesFromFile(file)

	for outCh != nil || errs != nil {
		select {
		case price, ok := <-outCh:
			if ok {
				fmt.Println(price)
			} else {
				outCh = nil
			}

		case err, ok := <-errs:
			if ok {
				fmt.Println("Error", err)
			} else {
				errs = nil
			}
		}
	}
}
