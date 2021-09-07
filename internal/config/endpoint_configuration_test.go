package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEndpointConfiguration_GetEndpointAddress(t *testing.T) {
	ec := EndpointConfiguration{
		Host: "0.0.0.0",
		Port: 777,
	}
	expected := "0.0.0.0:777"

	result := ec.GetEndpointAddress()

	assert.Equal(t, expected, result, "should return correct endpoint string")
}
