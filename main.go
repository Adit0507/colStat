package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
)

func run(filenames []string, op string, col int, out io.Writer) error {
	var opFunc statsFunc

	if len(filenames) == 0 {
		return ErrNoFiles
	}

	if col < 1 {
		return fmt.Errorf("%w: %d", ErrInvalidColumn, col)
	}

	switch op {
	case "sum":
		opFunc = sum
	case "avg":
		opFunc = avg

	default:
		return fmt.Errorf("%w: %s", ErrInvalidOperation, op)
	}

	// consolidate data extracted from given col. on each input file
	consolidate := make([]float64, 0)

	// channel to receive results or errors of operations
	resCh := make(chan []float64)
	errCh := make(chan error)
	doneCh := make(chan struct{})

	wg := sync.WaitGroup{}

	// loopin thrugh all files and create goroutine to process each one concurrently
	for _, fname := range filenames {
		wg.Add(1)
		go func(fname string) {
			defer wg.Done()

			f, err := os.Open(fname)
			if err != nil {
				errCh <- fmt.Errorf("cannot open file: %w", err)
				return
			}
			data, err := csv2float(f, col)
			if err != nil {
				errCh <- err
			}

			if err := f.Close(); err != nil {
				errCh <- err
			}

			resCh <- data
		}(fname)
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			return err
		case data := <-resCh:
			consolidate = append(consolidate, data...)
		case <-doneCh:
			_, err := fmt.Fprintln(out, opFunc(consolidate))
			return err
		}
	}
}

func main() {
	op := flag.String("op", "sum", "Operation to be executed")
	col := flag.Int("col", 1, "CSV column on which to execute operation")
	flag.Parse()

	if err := run(flag.Args(), *op, *col, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
