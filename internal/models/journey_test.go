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
		UserId:      1,
		JourneyId:   2,
		Address:     "Воронеж",
		Description: "Поездка на выходные",
	}
	journey.StartTime, _ = time.Parse(time.RFC822, startTimeStringRFC822)
	journey.EndTime, _ = time.Parse(time.RFC822, endTimeStringRFC822)

	result := journey.String()

	// "Journey: Id = 2, UserId = 1, Address = Воронеж, Description = Поездка на выходные, StartTime = 14 Aug 21 09:00 MST, EndTime = 15 Aug 21 21:00 MST"
	assert.Equal(t,
		fmt.Sprintf(
			"Journey: Id = %d, UserId = %d, Address = %s, Description = %s, StartTime = %s, EndTime = %s",
			journey.JourneyId,
			journey.UserId,
			journey.Address,
			journey.Description,
			startTimeStringRFC822,
			endTimeStringRFC822,
		), result)
}
