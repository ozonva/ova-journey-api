package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSwapMapKeyValue(t *testing.T) {
	var testTable = []struct {
		sourceMap map[int]string
		destMap   map[string]int
		err       error
	}{
		{
			sourceMap: map[int]string{
				1: "ova",
				2: "journey",
				3: "api",
			},
			destMap: map[string]int{
				"ova":     1,
				"journey": 2,
				"api":     3,
			},
			err: nil,
		},
		{
			sourceMap: map[int]string{
				1: "ova",
				2: "ova",
				3: "api",
			},
			destMap: nil,
			err:     ErrSourceMapMustBeUnique,
		},
		{
			sourceMap: nil,
			destMap:   nil,
			err:       nil,
		},
	}

	for _, testCase := range testTable {
		result, err := SwapMapKeyValue(testCase.sourceMap)

		assert.Equal(t, testCase.destMap, result)
		assert.Equal(t, testCase.err, err)
	}
}
