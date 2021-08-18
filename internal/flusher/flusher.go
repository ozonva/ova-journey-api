package flusher

import (
	"github.com/ozonva/ova-journey-api/internal/models"
	"github.com/ozonva/ova-journey-api/internal/repo"
	"github.com/ozonva/ova-journey-api/internal/utils"
)

type Flusher interface {
	Flush(journeys []models.Journey) []models.Journey
}

type flusher struct {
	chunkSize   int
	journeyRepo repo.Repo
}

func (f *flusher) Flush(journeys []models.Journey) []models.Journey {
	chunks, err := utils.SplitToChunks(journeys, f.chunkSize)
	if err != nil {
		return journeys
	}
	var failedJourneys []models.Journey = nil
	for i, chunk := range chunks {
		if err := f.journeyRepo.AddJourneys(chunk); err != nil {
			if failedJourneys == nil {
				failedJourneys = make([]models.Journey, 0, len(journeys)-i*f.chunkSize)
			}
			failedJourneys = append(failedJourneys, chunk...)
		}
	}
	return failedJourneys
}

func NewFlusher(chunkSize int, repo repo.Repo) Flusher {
	return &flusher{
		chunkSize:   chunkSize,
		journeyRepo: repo,
	}
}
