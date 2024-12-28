package repo

import (
	"context"

	model "github.com/pawelkuk/pictura-certamine/pkg/domain/user/model"
)

type Repo interface {
	Create(context.Context, *model.User) error
	Read(context.Context, *model.User) error
	Update(context.Context, *model.User) error
	Delete(context.Context, *model.User) error
	Query(context.Context, model.QueryFilter) ([]model.User, error)
}
