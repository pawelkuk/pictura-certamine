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
	ID             *string
	Email          *mail.Address
	FirstName      *string
	LastName       *string
	Birthdate      *time.Time
	PolicyAccepted *bool
}

type EntryQueryFilter struct {
	ID           *string
	ContestantID *string
	SessionID    *string
	Status       *EntryStatus
}
