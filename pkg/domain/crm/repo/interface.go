package repo

import (
	"context"

	"github.com/pawelkuk/pictura-certamine/pkg/domain/crm/model"
)

type Repo interface {
	Read(context.Context, *model.ContestantEntry) error
	Query(context.Context, model.ContestantEntryQueryFilter) ([]model.ContestantEntry, error)
}
