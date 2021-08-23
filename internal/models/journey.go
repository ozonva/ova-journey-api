package models

import (
	"fmt"
	"time"
)

//Journey - represents the journey description object
type Journey struct {
	JourneyID   uint64
	UserID      uint64
	Address     string
	Description string
	StartTime   time.Time
	EndTime     time.Time
}

func (j *Journey) String() string {
	return fmt.Sprintf(
		"Journey: Id = %d, UserID = %d, Address = %s, Description = %s, StartTime = %s, EndTime = %s",
		j.JourneyID,
		j.UserID,
		j.Address,
		j.Description,
		j.StartTime.Format(time.RFC822),
		j.EndTime.Format(time.RFC822),
	)
}

// NewJourney - creates new Journey object using arguments
func NewJourney(journeyID uint64, userId uint64, address string, description string, startTime time.Time, endTime time.Time) *Journey {
	return &Journey{
		JourneyID:   journeyID,
		UserID:      userId,
		Address:     address,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
	}
}
