package repo

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/ozonva/ova-journey-api/internal/models"
)

//Repo - represents the object for working with storage of Journeys
type Repo interface {
	AddJourney(ctx context.Context, journey models.Journey) (uint64, error)
	AddJourneysMulti(ctx context.Context, journeys []models.Journey) error
	ListJourneys(ctx context.Context, limit, offset uint64) ([]models.Journey, error)
	DescribeJourney(ctx context.Context, journeyID uint64) (*models.Journey, error)
	RemoveJourney(ctx context.Context, journeyID uint64) error
}

type repo struct {
	db *sqlx.DB
}

// NewRepo - creates new Journey repository using database
func NewRepo(db *sqlx.DB) Repo {
	return &repo{db: db}
}

func (r *repo) AddJourney(ctx context.Context, journey models.Journey) (uint64, error) {
	query := squirrel.
		Insert("journeys").
		Columns("user_id", "address", "description", "start_time", "end_time").
		Values(journey.UserID, journey.Address, journey.Description, journey.StartTime, journey.EndTime).
		Suffix("RETURNING \"journey_id\"").
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	var journeyID uint64
	err := query.QueryRowContext(ctx).Scan(&journeyID)
	if err != nil {
		return 0, err
	}
	return journeyID, nil
}

func (r *repo) AddJourneysMulti(ctx context.Context, journeys []models.Journey) error {
	query := squirrel.
		Insert("journeys").
		Columns("user_id", "address", "description", "start_time", "end_time").
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	for _, journey := range journeys {
		query = query.Values(journey.UserID, journey.Address, journey.Description, journey.StartTime, journey.EndTime)
	}

	_, err := query.ExecContext(ctx)
	return err

}

func (r *repo) ListJourneys(ctx context.Context, limit, offset uint64) ([]models.Journey, error) {
	query := squirrel.
		Select("journey_id", "user_id", "address", "description", "start_time", "end_time").
		From("journeys").
		Where(squirrel.Eq{"is_deleted": false}).
		Limit(limit).
		Offset(offset).
		OrderBy("journey_id ASC").
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var journeysList []models.Journey
	for rows.Next() {
		var journey models.Journey
		err = rows.Scan(
			&journey.JourneyID,
			&journey.UserID,
			&journey.Address,
			&journey.Description,
			&journey.StartTime,
			&journey.EndTime,
		)
		if err != nil {
			return journeysList, err
		}
		journeysList = append(journeysList, journey)
	}

	return journeysList, nil
}

func (r *repo) DescribeJourney(ctx context.Context, journeyID uint64) (*models.Journey, error) {
	query := squirrel.
		Select("journey_id", "user_id", "address", "description", "start_time", "end_time").
		From("journeys").
		Where(squirrel.And{squirrel.Eq{"journey_id": journeyID}, squirrel.Eq{"is_deleted": false}}).
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	var journey models.Journey
	err := query.QueryRowContext(ctx).
		Scan(
			&journey.JourneyID,
			&journey.UserID,
			&journey.Address,
			&journey.Description,
			&journey.StartTime,
			&journey.EndTime,
		)
	if err != nil {
		return nil, err
	}
	return &journey, nil
}

func (r *repo) RemoveJourney(ctx context.Context, journeyID uint64) error {
	query := squirrel.
		Update("journeys").
		Set("is_deleted", true).
		Where(squirrel.Eq{"journey_id": journeyID}).
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	_, err := query.ExecContext(ctx)
	return err
}
