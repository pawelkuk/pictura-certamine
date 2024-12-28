package repo

import (
	"context"

	model "github.com/pawelkuk/pictura-certamine/pkg/domain/user/model"
)

type FakeRepo struct {
	user    *model.User
	users   []model.User
	userErr error
}

func (r *FakeRepo) Create(_ context.Context, user *model.User) error {
	if r.userErr != nil {
		return r.userErr
	}
	*user = *r.user
	return nil
}
func (r *FakeRepo) Read(_ context.Context, user *model.User) error {
	if r.userErr != nil {
		return r.userErr
	}
	*user = *r.user
	return nil
}
func (r *FakeRepo) Update(_ context.Context, user *model.User) error {
	if r.userErr != nil {
		return r.userErr
	}
	*user = *r.user
	return nil
}
func (r *FakeRepo) Delete(_ context.Context, user *model.User) error {
	if r.userErr != nil {
		return r.userErr
	}
	*user = *r.user
	return nil
}

func (r *FakeRepo) Query(_ context.Context, queryFilter model.QueryFilter) ([]model.User, error) {
	if r.userErr != nil {
		return nil, r.userErr
	}
	return r.users, nil
}
