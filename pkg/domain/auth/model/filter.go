package model

// QueryFilter holds the available fields a query can be filtered on.
// We are using pointer because the With API mutates the value.
type QueryFilter struct {
	Token  *SessionToken
	UserID *int64
}
