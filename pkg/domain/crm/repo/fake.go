package repo

import (
	"context"

	"github.com/pawelkuk/pictura-certamine/pkg/domain/crm/model"
)

type FakeRepo struct {
	contestant    *model.ContestantEntry
	contestantErr error
}

func (r *FakeRepo) Read(_ context.Context, contestant *model.ContestantEntry) error {
	if r.contestantErr != nil {
		return r.contestantErr
	}
	*contestant = *r.contestant
	return nil
}

func (r *FakeRepo) Query(context.Context, model.ContestantEntryQueryFilter) ([]model.ContestantEntry, error) {
	return nil, nil
}
