package repo

import (
	"context"

	"github.com/pawelkuk/pictura-certamine/pkg/bus/contest/model"
)

type Repo interface {
	Create(context.Context, *model.Contestant) error
	Read(context.Context, *model.Contestant) error
	Update(context.Context, *model.Contestant) error
	Delete(context.Context, *model.Contestant) error
	Query(context.Context, model.ContestantQueryFilter) ([]model.Contestant, error)
}
