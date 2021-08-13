package utils

import "errors"

var (
	ErrSliceCannotBeNil        = errors.New("slice cannot be nil")
	ErrSliceCannotBeNilOrEmpty = errors.New("slice cannot be nil or empty")
	ErrIncorrectChunkSize      = errors.New("chunk size must be greater then zero")
	ErrSourceMapMustBeUnique   = errors.New("sourceMap values must be unique to swap")
)
