package model

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Token  SessionToken
	Expiry time.Time
	UserID int64
}

type SessionToken struct {
	Value string
}

func New(userid int64) Session {
	return Session{
		Token: SessionToken{
			Value: uuid.NewString(),
		},
		Expiry: time.Now().Add(24 * time.Hour),
		UserID: userid,
	}
}

func (s *Session) Refresh(duration time.Duration) {
	s.Expiry = time.Now().Add(duration)
}
