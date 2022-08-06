package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

type Options struct {
	MinValue  int
	MaxValue  int
	Arrays    int
	ArraySize int
}

const (
	MinArraySize = 2
)

func validateOptions(options Options) error {
	if options.MinValue < 0 {
		return fmt.Errorf("minimum value cannot be less than 0, got %d", options.MinValue)
	}

	if options.MaxValue < 0 {
		return fmt.Errorf("maximum value cannot be less than 0, got %d", options.MaxValue)
	}

	if options.MinValue >= options.MaxValue {
		return fmt.Errorf("minimum (%d) value cannot be greater than maximum (%d)",
			options.MinValue, options.MaxValue)
	}

	if options.Arrays <= 0 {
		return fmt.Errorf("number (%d) of arrays generated must be greater than 0", options.Arrays)
	}

	if options.ArraySize <= 0 {
		return fmt.Errorf("maximum array size (%d) generated must be greater than 0", options.ArraySize)
	}

	return nil
}

func getOptions() (Options, error) {
	var options Options

	flag.IntVar(&options.MinValue, "min", 0, "Minimum value to generate")
	flag.IntVar(&options.MaxValue, "max", 10_000, "Maximum value to generate")
	flag.IntVar(&options.Arrays, "arrays", rand.Intn(1_000)+1, "Number of arrays to generate")
	flag.IntVar(&options.ArraySize, "size", rand.Intn(100_000)+1, "Maximum array size to generate")

	flag.Parse()

	err := validateOptions(options)
	if err != nil {
		return Options{}, err
	}

	return options, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	options, err := getOptions()
	if err != nil {
		panic(err)
	}

	for i := 0; i < options.Arrays; i++ {
		arraySize := rand.Intn(options.ArraySize-MinArraySize) + MinArraySize

		for j := 0; j < arraySize; j++ {
			value := rand.Intn(options.MaxValue-options.MinValue) + options.MinValue
			fmt.Printf("%d ", value)
		}

		fmt.Printf("\n")
	}

}
