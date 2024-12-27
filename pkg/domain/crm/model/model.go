package model

import (
	"time"
)

type ContestantEntry struct {
	ID                string
	ContestantID      string
	Status            string
	ArtPieces         []ArtPiece
	Email             string
	PhoneNumber       string
	FirstName         string
	LastName          string
	ConsentConditions bool
	ConsentMarketing  bool
}

type ArtPiece struct {
	ID        int64
	Key       string
	CreatedAt time.Time
}
