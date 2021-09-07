package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrometheusConfiguration_GetEndpointAddress(t *testing.T) {
	pc := PrometheusConfiguration{
		Host: "0.0.0.0",
		Port: 777,
		Path: "/metrics",
	}
	expected := "0.0.0.0:777"

	result := pc.GetEndpointAddress()

	assert.Equal(t, expected, result, "should return correct endpoint string")
}
