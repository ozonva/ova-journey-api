package utils

// SwapMapKeyValue returns new map by swapping keys and values in original map
// If original map contains repeated values returns error
func SwapMapKeyValue(sourceMap map[int]string) (map[string]int, error) {
	if sourceMap == nil {
		return nil, nil
	}

	destMap := make(map[string]int, len(sourceMap))

	for key, value := range sourceMap {
		if _, found := destMap[value]; found {
			return nil, ErrSourceMapMustBeUnique
		}
		destMap[value] = key
	}

	return destMap, nil
}
