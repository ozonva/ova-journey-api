package utils

import (
	"github.com/ozonva/ova-journey-api/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testSplitToChunksResult struct {
	chunks [][]models.Journey
	err    error
}

type testTestSliceToMapResult struct {
	mapDst map[uint64]models.Journey
	err    error
}

func TestSplitToChunks(t *testing.T) {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	testTable := []struct {
		slice     []models.Journey
		chunkSize int
		result    testSplitToChunksResult
	}{
		{
			slice: []models.Journey{
				*models.NewJourney(1, 1, "Воронеж", "", now, tomorrow),
				*models.NewJourney(2, 2, "Уфа", "", now, tomorrow),
				*models.NewJourney(3, 3, "Москва", "Командировка", now, tomorrow),
				*models.NewJourney(4, 4, "Хабаровск", "", now, tomorrow),
				*models.NewJourney(5, 5, "Лондон", "Командировка", now, tomorrow),
			},
			chunkSize: 2,
			result: testSplitToChunksResult{
				chunks: [][]models.Journey{
					{
						*models.NewJourney(1, 1, "Воронеж", "", now, tomorrow),
						*models.NewJourney(2, 2, "Уфа", "", now, tomorrow),
					},
					{
						*models.NewJourney(3, 3, "Москва", "Командировка", now, tomorrow),
						*models.NewJourney(4, 4, "Хабаровск", "", now, tomorrow),
					},
					{
						*models.NewJourney(5, 5, "Лондон", "Командировка", now, tomorrow),
					},
				},
				err: nil,
			},
		},
		{
			slice:     nil,
			chunkSize: 2,
			result: testSplitToChunksResult{
				chunks: nil,
				err:    ErrSliceCannotBeNil,
			},
		},
		{
			slice:     []models.Journey{},
			chunkSize: 0,
			result: testSplitToChunksResult{
				chunks: nil,
				err:    ErrIncorrectChunkSize,
			},
		},
	}

	for _, testCase := range testTable {
		result, err := SplitToChunks(testCase.slice, testCase.chunkSize)
		assert.Equal(t, testCase.result.chunks, result)
		assert.Equal(t, testCase.result.err, err)
	}
}

func TestSliceToMap(t *testing.T) {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	testTable := []struct {
		slice  []models.Journey
		result testTestSliceToMapResult
	}{
		{
			slice: []models.Journey{
				*models.NewJourney(1, 1, "Воронеж", "", now, tomorrow),
				*models.NewJourney(2, 2, "Уфа", "", now, tomorrow),
				*models.NewJourney(3, 3, "Москва", "Командировка", now, tomorrow),
			},
			result: testTestSliceToMapResult{
				mapDst: map[uint64]models.Journey{
					1: *models.NewJourney(1, 1, "Воронеж", "", now, tomorrow),
					2: *models.NewJourney(2, 2, "Уфа", "", now, tomorrow),
					3: *models.NewJourney(3, 3, "Москва", "Командировка", now, tomorrow),
				},
				err: nil,
			},
		},
		{
			slice: nil,
			result: testTestSliceToMapResult{
				mapDst: nil,
				err:    ErrSliceCannotBeNilOrEmpty,
			},
		},
		{
			slice: []models.Journey{},
			result: testTestSliceToMapResult{
				mapDst: nil,
				err:    ErrSliceCannotBeNilOrEmpty,
			},
		},
	}

	for _, testCase := range testTable {
		mapDst, err := SliceToMap(testCase.slice)
		assert.Equal(t, testCase.result.mapDst, mapDst)
		assert.Equal(t, testCase.result.err, err)
	}
}
