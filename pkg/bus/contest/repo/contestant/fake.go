package repo

import (
	"context"

	"github.com/pawelkuk/pictura-certamine/pkg/bus/contest/model"
)

type FakeRepo struct {
	contestant    *model.Contestant
	contestantErr error
}

func (r *FakeRepo) Create(_ context.Context, contestant *model.Contestant) error {
	if r.contestantErr != nil {
		return r.contestantErr
	}
	*contestant = *r.contestant
	return nil
}
func (r *FakeRepo) Read(_ context.Context, contestant *model.Contestant) error {
	if r.contestantErr != nil {
		return r.contestantErr
	}
	*contestant = *r.contestant
	return nil
}
func (r *FakeRepo) Update(_ context.Context, contestant *model.Contestant) error {
	if r.contestantErr != nil {
		return r.contestantErr
	}
	*contestant = *r.contestant
	return nil
}
func (r *FakeRepo) Delete(_ context.Context, contestant *model.Contestant) error {
	if r.contestantErr != nil {
		return r.contestantErr
	}
	*contestant = *r.contestant
	return nil
}
