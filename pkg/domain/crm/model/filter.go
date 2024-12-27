package model

import (
	"time"
)

type ContestantEntryQueryFilter struct {
	ID                *string
	Email             *string
	FirstName         *string
	LastName          *string
	ConsentConditions *bool
	ConsentMarketing  *bool
	ContestantID      *string
	Status            *string
	Token             *string
	TokenExpiry       *time.Time
}
