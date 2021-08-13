package models

import (
	"fmt"
	"time"
)

type Journey struct {
	JourneyId   uint64
	UserId      uint64
	Address     string
	Description string
	StartTime   time.Time
	EndTime     time.Time
}

func (j *Journey) String() string {
	return fmt.Sprintf(
		"Journey: Id = %d, UserId = %d, Address = %s, Description = %s, StartTime = %s, EndTime = %s",
		j.JourneyId,
		j.UserId,
		j.Address,
		j.Description,
		j.StartTime.Format(time.RFC822),
		j.EndTime.Format(time.RFC822),
	)
}

func NewJourney(journeyId uint64, userId uint64, address string, description string, startTime time.Time, endTime time.Time) *Journey {
	return &Journey{
		JourneyId:   journeyId,
		UserId:      userId,
		Address:     address,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
	}
}
