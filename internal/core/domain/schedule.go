package domain

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	MedicationName string    `json:"meducation_name"`
	Frequency      int       `json:"frequency"`
	Duration       int       `json:"duration"`
	IsPermanent    bool      `json:"is_permanent"`
	StartDate      time.Time `json:"start_date"`
}

func NewSchedule(userID string, medicationName string, frequency int, duration int) *Schedule {
	if frequency <= 0 {
		frequency = 1
	}
	if duration < 0 {
		duration = 0
	}

	isPermanent := duration == 0

	return &Schedule{
		ID:             uuid.NewString(),
		UserID:         userID,
		MedicationName: medicationName,
		Frequency:      frequency,
		Duration:       duration,
		IsPermanent:    isPermanent,
		StartDate:      time.Now(),
	}
}
