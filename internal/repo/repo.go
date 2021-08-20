package repo

import "github.com/ozonva/ova-journey-api/internal/models"

type Repo interface {
	AddJourneys(journeys []models.Journey) error
	ListJourneys(limit, offset uint64) ([]models.Journey, error)
	DescribeJourney(journeyID uint64) (*models.Journey, error)
}
