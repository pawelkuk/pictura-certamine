package repo

import (
	"context"

	"github.com/pawelkuk/pictura-certamine/pkg/domain/contest/model"
)

type Repo interface {
	Create(context.Context, *model.Entry) error
	Read(context.Context, *model.Entry) error
	Update(context.Context, *model.Entry) error
	Delete(context.Context, *model.Entry) error
	Query(context.Context, model.EntryQueryFilter) ([]model.Entry, error)
}
