package flusher

import (
	"context"

	"github.com/ozonva/ova-journey-api/internal/models"
	"github.com/ozonva/ova-journey-api/internal/repo"
	"github.com/ozonva/ova-journey-api/internal/utils"
)

// Flusher represents the object used for flushing journey to data storage
type Flusher interface {
	// Flush - flush journeys to the storage and returns journeys slice that was not saved
	Flush(ctx context.Context, journeys []models.Journey) []models.Journey
}

type flusher struct {
	chunkSize   int
	journeyRepo repo.Repo
}

// Flush - flush journeys to the repo.Repo and returns journeys slice that was not saved
func (f *flusher) Flush(ctx context.Context, journeys []models.Journey) []models.Journey {
	chunks, err := utils.SplitToChunks(journeys, f.chunkSize)
	if err != nil {
		return journeys
	}
	var failedJourneys []models.Journey
	for i, chunk := range chunks {
		if err := f.journeyRepo.AddJourneysMulti(ctx, chunk); err != nil {
			if failedJourneys == nil {
				failedJourneys = make([]models.Journey, 0, len(journeys)-i*f.chunkSize)
			}
			failedJourneys = append(failedJourneys, chunk...)
		}
	}
	return failedJourneys
}

// NewFlusher return Flusher for saving journeys to repo.Repo with splitting on chunkSize batches.
func NewFlusher(chunkSize int, repo repo.Repo) Flusher {
	return &flusher{
		chunkSize:   chunkSize,
		journeyRepo: repo,
	}
}
