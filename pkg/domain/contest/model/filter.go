package model

import (
	"net/mail"
	"time"
)

type ContestQueryFilter struct {
	ID       *string
	Name     *string
	Slug     *Slug
	Start    *time.Time
	End      *time.Time
	IsActive *bool
}

type ContestantQueryFilter struct {
	ID                *string
	Email             *mail.Address
	FirstName         *string
	LastName          *string
	ConsentConditions *bool
	ConsentMarketing  *bool
}

type EntryQueryFilter struct {
	ID           *string
	ContestantID *string
	Status       *EntryStatus
}
