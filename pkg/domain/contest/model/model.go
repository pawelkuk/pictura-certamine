package model

import (
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/gosimple/slug"
	"github.com/hashicorp/go-multierror"
	"github.com/nyaruka/phonenumbers"
)

type EntryStatus string
type ParseError struct {
	Field string
	Err   error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s: %e", e.Field, e.Err)
}

var (
	EntryStatusPending               EntryStatus = "Pending"
	EntryStatusSubmitted             EntryStatus = "Submitted"
	EntryStatusConfirmationEmailSent EntryStatus = "ConfirmationEmailSent"
	EntryStatusConfirmed             EntryStatus = "Confirmed"
)

var validEntryStatusMap = map[string]EntryStatus{
	"Pending":               EntryStatusPending,
	"Submitted":             EntryStatusSubmitted,
	"ConfirmationEmailSent": EntryStatusConfirmationEmailSent,
	"Confirmed":             EntryStatusConfirmed,
}

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
	PhoneNumber    *phonenumbers.PhoneNumber
	FirstName      string
	LastName       string
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

func ParseSlug(raw string) Slug {
	val := slug.Make(raw)
	return Slug{
		Value: val,
	}
}

func ParseStatus(status string) (EntryStatus, error) {
	s, ok := validEntryStatusMap[status]
	if !ok {
		return "", fmt.Errorf("invalid status: %s", status)
	}
	return s, nil
}

func ParseContestant(
	id string,
	email string,
	phoneNumber string,
	firstName string,
	lastName string,
	birthdate string,
	policyAccepted string,
) (*Contestant, error) {
	var errs *multierror.Error
	if id == "" {
		id = generateID()
	}
	a, err := mail.ParseAddress(email)
	if err != nil {
		errs = multierror.Append(errs, &ParseError{Field: "Email", Err: err})
	}
	pn, err := phonenumbers.Parse(phoneNumber, "us")
	if err != nil {
		errs = multierror.Append(errs, &ParseError{Field: "PhoneNumber", Err: err})
	}
	if firstName == "" {
		errs = multierror.Append(errs, &ParseError{Field: "FirstName", Err: errors.New("first name can't be empty")})
	}
	if lastName == "" {
		errs = multierror.Append(errs, &ParseError{Field: "LastName", Err: errors.New("last name can't be empty")})
	}
	b, err := time.Parse(time.DateOnly, birthdate)
	if err != nil {
		errs = multierror.Append(errs, &ParseError{Field: "Birthdate", Err: err})
	}
	var pa bool
	if policyAccepted != "" && policyAccepted != "no" {
		pa = true
	}

	if errs.Len() == 0 {
		return &Contestant{
			ID:             id,
			Email:          *a,
			PhoneNumber:    pn,
			FirstName:      firstName,
			LastName:       lastName,
			Birthdate:      b,
			PolicyAccepted: pa,
		}, nil
	} else {
		return nil, errs
	}
}

func generateID() string {
	return "abcd"
}
