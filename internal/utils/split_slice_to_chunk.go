package utils

import "errors"

// SplitSliceToChunk - split slice to chunks of equal size (chunkSize), except last chunk that contains last elements of slice
func SplitSliceToChunk(slice []int, chunkSize int) ([][]int, error) {
	if slice == nil {
		return nil, errors.New("slice cannot be nil")
	}

	if chunkSize < 1 {
		return nil, errors.New("chunk size must be greater then zero")
	}

	var chunksCount = len(slice) / chunkSize
	if len(slice)%chunkSize > 0 {
		chunksCount++
	}

	var chunks = make([][]int, chunksCount, chunksCount)
	var start, end = 0, 0
	for i := 0; i < chunksCount; i++ {
		start = end
		end += chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks[i] = slice[start:end]
	}
	return chunks, nil
}
