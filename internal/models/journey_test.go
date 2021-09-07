package models

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJourney_String(t *testing.T) {
	startTimeStringRFC822 := "14 Aug 21 09:00 MST"
	endTimeStringRFC822 := "15 Aug 21 21:00 MST"
	journey := Journey{
		UserID:      1,
		JourneyID:   2,
		Address:     "Воронеж",
		Description: "Поездка на выходные",
	}
	journey.StartTime, _ = time.Parse(time.RFC822, startTimeStringRFC822)
	journey.EndTime, _ = time.Parse(time.RFC822, endTimeStringRFC822)

	result := journey.String()

	// "Journey: Id = 2, UserID = 1, Address = Воронеж, Description = Поездка на выходные, StartTime = 14 Aug 21 09:00 MST, EndTime = 15 Aug 21 21:00 MST"
	assert.Equal(t,
		fmt.Sprintf(
			"Journey: Id = %d, UserID = %d, Address = %s, Description = %s, StartTime = %s, EndTime = %s",
			journey.JourneyID,
			journey.UserID,
			journey.Address,
			journey.Description,
			startTimeStringRFC822,
			endTimeStringRFC822,
		), result)
}

func TestNewJourney(t *testing.T) {
	journey := &Journey{
		JourneyID:   2,
		UserID:      1,
		Address:     "Воронеж",
		Description: "Поездка на выходные",
		StartTime:   time.Now(),
		EndTime:     time.Now(),
	}

	result := NewJourney(
		journey.JourneyID,
		journey.UserID,
		journey.Address,
		journey.Description,
		journey.StartTime,
		journey.EndTime,
	)

	assert.Equal(t, journey, result, "should return equal journey")
}
