package utils

import "errors"

// SwapMapKeyValue returns new map by swapping keys and values in original map
// If original map contains repeated values returns error
func SwapMapKeyValue(sourceMap map[int]string) (map[string]int, error) {
	if sourceMap == nil {
		return nil, nil
	}

	var destMap map[string]int = make(map[string]int, len(sourceMap))

	for key, value := range sourceMap {
		if _, found := destMap[value]; found {
			return nil, errors.New("sourceMap values must be unique to swap")
		}
		destMap[value] = key
	}

	return destMap, nil
}
