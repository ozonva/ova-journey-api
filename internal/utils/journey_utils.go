package utils

import (
	"github.com/ozonva/ova-journey-api/internal/models"
)

// SplitToChunks - split models.Journey slice into batches with size of chinkSize
func SplitToChunks(slice []models.Journey, chunkSize int) ([][]models.Journey, error) {
	if slice == nil {
		return nil, ErrSliceCannotBeNil
	}

	if chunkSize < 1 {
		return nil, ErrIncorrectChunkSize
	}

	sliceLength := len(slice)
	chunksCount := sliceLength / chunkSize
	if sliceLength%chunkSize > 0 {
		chunksCount++
	}

	var chunks = make([][]models.Journey, chunksCount)

	for i, end := 0, 0; i < chunksCount; i++ {
		start := end
		end += chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks[i] = slice[start:end]
	}
	return chunks, nil
}

//SliceToMap - create map with JourneyID as key and Journey as value from models.Journey slice
func SliceToMap(srcSlice []models.Journey) (map[uint64]models.Journey, error) {
	if len(srcSlice) == 0 {
		return nil, ErrSliceCannotBeNilOrEmpty
	}

	destMap := make(map[uint64]models.Journey, len(srcSlice))

	for _, journey := range srcSlice {
		destMap[journey.JourneyID] = journey
	}

	return destMap, nil
}
