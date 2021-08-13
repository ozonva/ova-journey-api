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
	}{
		{
			slice:     []int{1, 2},
			chunkSize: 1,
			result:    [][]int{[]int{1}, []int{2}},
		},
		{
			slice:     []int{1, 2, 3, 4, 5},
			chunkSize: 2,
			result:    [][]int{[]int{1, 2}, []int{3, 4}, []int{5}},
		},
		{
			slice:     nil,
			chunkSize: 2,
			result:    nil,
		},
		{
			slice:     []int{1, 2, 3, 4, 5},
			chunkSize: 0,
			result:    nil,
		},
		{
			slice:     []int{},
			chunkSize: 1,
			result:    [][]int{},
		},
		{
			slice:     []int{1, 2},
			chunkSize: 10,
			result:    [][]int{[]int{1, 2}},
		},
	}

	for _, testCase := range testTable {
		result, err := SplitSliceToChunk(testCase.slice, testCase.chunkSize)

		assert.Equal(t, testCase.result, result)
		assert.Equal(t, testCase.err, err)

	}
}
