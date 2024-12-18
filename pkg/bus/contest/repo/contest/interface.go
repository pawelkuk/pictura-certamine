package repo

import (
	"context"

	"github.com/pawelkuk/pictura-certamine/pkg/bus/contest/model"
)

type Repo interface {
	Create(context.Context, *model.Contest) error
	Read(context.Context, *model.Contest) error
	Update(context.Context, *model.Contest) error
	Delete(context.Context, *model.Contest) error
	Query(context.Context, model.QueryFilter) ([]model.Contest, error)
}
