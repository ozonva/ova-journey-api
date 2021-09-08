package repo

import (
	"context"
	"github.com/ozonva/ova-journey-api/internal/models"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

var db *sqlx.DB
var repository Repo
var journeysTable []models.Journey = []models.Journey{
	{JourneyID: 1, UserID: 1, Address: "Уфа", Description: ""},
	{JourneyID: 2, UserID: 2, Address: "Москва", Description: ""},
	{JourneyID: 3, UserID: 2, Address: "Лондон", Description: ""},
	{JourneyID: 4, UserID: 3, Address: "Новосибирск", Description: ""},
}

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to docker")
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "13",
			Env: []string{
				"POSTGRES_USER=user_name",
				"POSTGRES_PASSWORD=secret",
				"POSTGRES_DB=dbname",
				"listen_addresses = '*'",
			},
			ExposedPorts: []string{"5432"},
			PortBindings: map[docker.Port][]docker.PortBinding{
				"5432": {
					{HostIP: "0.0.0.0", HostPort: "5432"},
				},
			},
		}, func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}

		})
	//defer resource.Close()

	if err != nil {
		log.Fatal().Err(err).Msg("Could not start resource")
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		//sqlx.Open(configuration.Driver, configuration.GetDataSourceName())
		db, err = sqlx.Open(
			"pgx",
			"postgres://user_name:secret@localhost:5432/dbname?sslmode=disable",
		)

		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		log.Fatal().Err(err).Msg("Could not connect to docker")
	}

	err = goose.Run("up", db.DB, "../../migrations")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not init migrations")
	}

	repository = NewRepo(db)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatal().Err(err).Msg("Could not purge resource")
	}

	os.Exit(code)
}

func TestRepo_AddJourney(t *testing.T) {
	id, err := repository.AddJourney(context.Background(), journeysTable[0])

	assert.NoError(t, err)
	assert.Greater(t, id, uint64(0), "Id must be greater then 0")
}

func TestRepo_DescribeJourney(t *testing.T) {
	journey, err := repository.DescribeJourney(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, journey)
}

func TestRepo_ListJourneys(t *testing.T) {
	journeys, err := repository.ListJourneys(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, journeys)
}

func TestRepo_MultiAddJourneys(t *testing.T) {
	journeysIDs, err := repository.MultiAddJourneys(context.Background(), journeysTable)

	assert.NoError(t, err)
	assert.NotNil(t, journeysIDs)
}

func TestRepo_UpdateJourney(t *testing.T) {
	testJourney := journeysTable[0]
	id, _ := repository.AddJourney(context.Background(), journeysTable[0])
	testJourney.JourneyID = id
	testJourney.Address = "changedAddress"

	err := repository.UpdateJourney(context.Background(), testJourney)
	assert.NoError(t, err)
}

func TestRepo_RemoveJourney(t *testing.T) {
	id, _ := repository.AddJourney(context.Background(), journeysTable[0])

	err := repository.RemoveJourney(context.Background(), id)
	assert.NoError(t, err)
}
