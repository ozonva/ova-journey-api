package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitSliceToChunk(t *testing.T) {
	var testTable = []struct {
		slice     []int
		chunkSize int
		result    [][]int
		err       error
	}{
		{
			slice:     []int{1, 2},
			chunkSize: 1,
			result:    [][]int{{1}, {2}},
			err:       nil,
		},
		{
			slice:     []int{1, 2, 3, 4, 5},
			chunkSize: 2,
			result:    [][]int{{1, 2}, {3, 4}, {5}},
			err:       nil,
		},
		{
			slice:     nil,
			chunkSize: 2,
			result:    nil,
			err:       ErrSliceCannotBeNil,
		},
		{
			slice:     []int{1, 2, 3, 4, 5},
			chunkSize: 0,
			result:    nil,
			err:       ErrIncorrectChunkSize,
		},
		{
			slice:     []int{},
			chunkSize: 1,
			result:    [][]int{},
			err:       nil,
		},
		{
			slice:     []int{1, 2},
			chunkSize: 10,
			result:    [][]int{{1, 2}},
			err:       nil,
		},
	}

	for _, testCase := range testTable {
		result, err := SplitSliceToChunk(testCase.slice, testCase.chunkSize)
		assert.Equal(t, testCase.result, result)
		assert.Equal(t, testCase.err, err)
	}
}
