package utils

// FilterSlice returns new slice with elements from original slice not presented in exclusion slice
func FilterSlice(slice []string, excluded []string) []string {
	if slice == nil {
		return nil
	}

	var filteredSlice = make([]string, 0, len(slice))

	if len(slice) > 0 {
		// prepare filtration map
		var excludedMap = make(map[string]struct{}, len(excluded))
		for _, value := range excluded {
			excludedMap[value] = struct{}{}
		}
		// filter slice by map
		for _, value := range slice {
			if _, found := excludedMap[value]; !found {
				filteredSlice = append(filteredSlice, value)
			}
		}
	}

	return filteredSlice
}
