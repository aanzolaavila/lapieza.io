package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getPrices(file *os.File) ([]int64, error) {
	scanner := bufio.NewScanner(file)
	var prices []int64

	for scanner.Scan() {
		line := scanner.Text()

		strs := strings.Split(line, " ")

		for _, strNumber := range strs {
			num, err := strconv.ParseInt(strNumber, 10, 64)

			if err != nil {
				return nil, err
			}

			prices = append(prices, num)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return prices, nil
}

func main() {
	var fileName string
	flag.StringVar(&fileName, "file", "", "Specify filename. Default is STDIN.")
	flag.Parse()

	var file *os.File
	if fileName == "" {
		file = os.Stdin
	} else {
		var err error
		file, err = os.Open(fileName)
		if err != nil {
			fmt.Errorf("Failed to open %s: %v", fileName, err)
			panic(err)
		}
	}

	prices, err := getPrices(file)
	if err != nil {
		fmt.Errorf("Failed to read %s: %v", fileName, err)
		panic(err)
	}

	fmt.Println(prices)
}
