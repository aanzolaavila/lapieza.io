package prices

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aanzolaavila/lapieza.io/internal/types"
)

func GetInputFile() (file *os.File, err error) {
	file = os.Stdin

	var fileName string
	flag.StringVar(&fileName, "file", "", "Specify filename. Default is STDIN.")
	flag.Parse()

	if fileName != "" {
		file, err = os.Open(fileName)
		if err != nil {
			err = fmt.Errorf("failed to open %s: %v", fileName, err)
			return nil, err
		}
	}

	return file, nil
}

func getPricesFromLine(line string) ([]types.Price, error) {
	var prices []types.Price

	strs := strings.Split(line, " ")
	for _, strNumber := range strs {
		num, err := strconv.ParseInt(strNumber, 10, 64)

		if err != nil {
			return nil, err
		}

		prices = append(prices, types.Price(num))
	}

	return prices, nil
}

func streamPrices(file *os.File, inCh chan<- []types.Price, errs chan<- error) {
	scanner := bufio.NewScanner(file)

	defer close(inCh)
	defer close(errs)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			break
		}

		prices, err := getPricesFromLine(line)
		if err == nil {
			inCh <- prices
		} else {
			errs <- err
		}
	}

	if err := scanner.Err(); err != nil {
		err = fmt.Errorf("failed to read %s: %v", file.Name(), err)
		errs <- err
	}
}

func calculateMaxEarnings(prices []types.Price) types.Price {
	var maxDiff types.Price = 0

	for i := 0; i < len(prices); i++ {
		for j := i + 1; j < len(prices); j++ {
			start := prices[i]
			end := prices[j]

			diff := end - start
			if diff > maxDiff {
				maxDiff = diff
			}
		}
	}

	return maxDiff
}

func processPrices(inCh <-chan []types.Price, outCh chan<- types.Price) {
	for prices := range inCh {
		outCh <- calculateMaxEarnings(prices)
	}
}

func ProcessPricesFromFile(file *os.File) (<-chan types.Price, <-chan error) {
	inCh := make(chan []types.Price)
	outCh := make(chan types.Price)
	errs := make(chan error)

	go func() {
		streamPrices(file, inCh, errs)
	}()

	go func() {
		processPrices(inCh, outCh)
	}()

	return outCh, errs
}
