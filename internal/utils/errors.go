package utils

import "errors"

var (
	// ErrSliceCannotBeNil - occurs when slice used as function argument is nil
	ErrSliceCannotBeNil = errors.New("slice cannot be nil")

	// ErrSliceCannotBeNilOrEmpty - occurs when slice used as function argument is nil or empty
	ErrSliceCannotBeNilOrEmpty = errors.New("slice cannot be nil or empty")

	// ErrIncorrectChunkSize - occurs when chunkSize is incorrect.
	ErrIncorrectChunkSize = errors.New("chunk size must be greater then zero")

	//ErrSourceMapMustBeUnique - occurs when map used for swap keys and values contains non-unique keys
	ErrSourceMapMustBeUnique = errors.New("sourceMap values must be unique to swap")
)
