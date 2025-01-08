package model

import (
	"net/mail"
)

// QueryFilter holds the available fields a query can be filtered on.
// We are using pointer because the With API mutates the value.
type QueryFilter struct {
	ID                 *int64
	Email              *mail.Address
	AuthorizationToken *string
	ActivationToken    *string
	PasswordResetToken *string
}
