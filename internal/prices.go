package prices

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

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
		if strNumber == "" {
			continue
		}

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

	const maxCapacity int = 10 * 1024 * 1024 // 10 MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

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
	if len(prices) == 0 {
		return 0
	}

	var maxDiff types.Price = 0
	var minValue types.Price = prices[0]

	for i := 0; i < len(prices); i++ {
		if minValue > prices[i] {
			minValue = prices[i]
		}

		diff := prices[i] - minValue
		if diff > maxDiff {
			maxDiff = diff
		}
	}

	return maxDiff
}

func processPrices(inCh <-chan []types.Price, outCh chan<- types.Price) {
	for prices := range inCh {
		outCh <- calculateMaxEarnings(prices)
	}
}

func processPricesConcurrent(inCh <-chan []types.Price, outCh chan<- types.Price) {
	const goroutines = 4

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			processPrices(inCh, outCh)
			wg.Done()
		}()
	}

	wg.Wait()
	close(outCh)
}

func ProcessPricesFromFile(file *os.File) (<-chan types.Price, <-chan error) {
	inCh := make(chan []types.Price, 5)
	outCh := make(chan types.Price, 5)
	errs := make(chan error, 5)

	go func() {
		streamPrices(file, inCh, errs)
	}()

	go func() {
		processPricesConcurrent(inCh, outCh)
	}()

	return outCh, errs
}
