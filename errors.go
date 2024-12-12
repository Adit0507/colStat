package main

import "errors"

var (
	ErrNotNumber = errors.New("data is not numeric")
	ErrInvalidColumn = errors.New("invalid column no.")
	ErrNoFiles = errors.New("no input files")
	ErrInvalidOperation = errors.New("invalid operation")
)