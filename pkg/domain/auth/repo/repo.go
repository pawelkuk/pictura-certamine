package repo

import (
	"context"

	model "github.com/pawelkuk/pictura-certamine/pkg/domain/auth/model"
)

type Repo interface {
	Create(context.Context, *model.Session) error
	Read(context.Context, *model.Session) error
	Update(context.Context, *model.Session) error
	Delete(context.Context, *model.Session) error
	Query(context.Context, model.QueryFilter) ([]model.Session, error)
}
