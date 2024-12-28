package repo

import (
	"context"

	model "github.com/pawelkuk/pictura-certamine/pkg/domain/auth/model"
)

type FakeRepo struct {
	session    *model.Session
	sessionErr error
}

func (r *FakeRepo) Create(_ context.Context, task *model.Session) error {
	if r.sessionErr != nil {
		return r.sessionErr
	}
	*task = *r.session
	return nil
}
func (r *FakeRepo) Read(_ context.Context, task *model.Session) error {
	if r.sessionErr != nil {
		return r.sessionErr
	}
	*task = *r.session
	return nil
}
func (r *FakeRepo) Update(_ context.Context, task *model.Session) error {
	if r.sessionErr != nil {
		return r.sessionErr
	}
	*task = *r.session
	return nil
}
func (r *FakeRepo) Delete(_ context.Context, task *model.Session) error {
	if r.sessionErr != nil {
		return r.sessionErr
	}
	*task = *r.session
	return nil
}
