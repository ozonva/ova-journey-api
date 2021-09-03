package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatabaseConfiguration_GetDataSourceName(t *testing.T) {
	dbConfig := DatabaseConfiguration{
		Host:     "0.0.0.0",
		Port:     777,
		User:     "tu",
		Password: "tp",
		Name:     "tn",
		SslMode:  "disabled",
		Driver:   "pgx",
	}

	expected := "host=0.0.0.0 port=777 dbname=tn user=tu password=tp sslmode=disabled"

	result := dbConfig.GetDataSourceName()

	assert.Equal(t, expected, result, "should return correct dsn string")

}
