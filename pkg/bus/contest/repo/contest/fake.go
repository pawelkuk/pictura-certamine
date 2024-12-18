package repo

import (
	"context"

	"github.com/pawelkuk/pictura-certamine/pkg/bus/contest/model"
)

type FakeRepo struct {
	contest    *model.Contest
	contestErr error
}

func (r *FakeRepo) Create(_ context.Context, contest *model.Contest) error {
	if r.contestErr != nil {
		return r.contestErr
	}
	*contest = *r.contest
	return nil
}
func (r *FakeRepo) Read(_ context.Context, contest *model.Contest) error {
	if r.contestErr != nil {
		return r.contestErr
	}
	*contest = *r.contest
	return nil
}
func (r *FakeRepo) Update(_ context.Context, contest *model.Contest) error {
	if r.contestErr != nil {
		return r.contestErr
	}
	*contest = *r.contest
	return nil
}
func (r *FakeRepo) Delete(_ context.Context, contest *model.Contest) error {
	if r.contestErr != nil {
		return r.contestErr
	}
	*contest = *r.contest
	return nil
}
