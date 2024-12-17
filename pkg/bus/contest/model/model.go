package model

import (
	"net/mail"
	"time"
)

type EntryStatus string

var (
	EntryStatusPending               EntryStatus = "Pending"
	EntryStatusSubmitted             EntryStatus = "Submitted"
	EntryStatusConfirmationEmailSent EntryStatus = "ConfirmationEmailSent"
	EntryStatusConfirmed             EntryStatus = "Confirmed"
)

type Entry struct {
	ID           string
	ContestantID string
	SessionID    string
	Status       EntryStatus
	ArtPieces    []ArtPiece
}

type ArtPiece struct {
	ID        int64
	Key       string
	CreatedAt time.Time
}

type Contestant struct {
	ID             string
	Email          mail.Address
	FirstName      string
	Surname        string
	Birthdate      time.Time
	PolicyAccepted bool
}

type Contest struct {
	ID       string
	Name     string
	Slug     Slug
	Start    time.Time
	End      time.Time
	IsActive bool
}

type Slug struct {
	Value string
}
