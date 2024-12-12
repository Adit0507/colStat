package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

func sum(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}

	return sum
}

func avg(data []float64) float64 {
	return sum(data) / float64(len(data))
}

type statsFunc func(data []float64) float64

func csv2float(r io.Reader, col int) ([]float64, error) {
	// csv reaer used to read in data from csv file
	cr := csv.NewReader(r)

	col--
	// read all csv data
	allData, err := cr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Cannot read data from the file: %w", err)
	}

	var data []float64
	for i, row := range allData {
		if i == 0{
			continue
		}

		if len(row) <= col {
			return nil, fmt.Errorf("%w: File has only %d columns", ErrInvalidColumn, len(row))
		}

		v, err := strconv.ParseFloat(row[col], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, err)
		}

		data = append(data, v)
	}

	return data, nil
}