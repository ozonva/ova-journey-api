package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilterSlice(t *testing.T) {
	var testTable = []struct {
		slice    []string
		excluded []string
		result   []string
	}{
		{
			slice:    []string{"ova", "journey", "api"},
			excluded: []string{"journey"},
			result:   []string{"ova", "api"},
		},
		{
			slice:    []string{"ova", "journey", "api"},
			excluded: []string{"another"},
			result:   []string{"ova", "journey", "api"},
		},
		{
			slice:    []string{"ova", "journey", "api"},
			excluded: nil,
			result:   []string{"ova", "journey", "api"},
		},
		{
			slice:    nil,
			excluded: nil,
			result:   nil,
		},
		{
			slice:    []string{},
			excluded: nil,
			result:   []string{},
		},
	}

	for _, testCase := range testTable {
		result := FilterSlice(testCase.slice, testCase.excluded)

		assert.Equal(t, testCase.result, result)
	}
}
