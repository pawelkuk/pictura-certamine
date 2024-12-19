package repo

import (
	"context"

	"github.com/pawelkuk/pictura-certamine/pkg/domain/contest/model"
)

type FakeRepo struct {
	entry    *model.Entry
	entryErr error
}

func (r *FakeRepo) Create(_ context.Context, entry *model.Entry) error {
	if r.entryErr != nil {
		return r.entryErr
	}
	*entry = *r.entry
	return nil
}
func (r *FakeRepo) Read(_ context.Context, entry *model.Entry) error {
	if r.entryErr != nil {
		return r.entryErr
	}
	*entry = *r.entry
	return nil
}
func (r *FakeRepo) Update(_ context.Context, entry *model.Entry) error {
	if r.entryErr != nil {
		return r.entryErr
	}
	*entry = *r.entry
	return nil
}
func (r *FakeRepo) Delete(_ context.Context, entry *model.Entry) error {
	if r.entryErr != nil {
		return r.entryErr
	}
	*entry = *r.entry
	return nil
}
