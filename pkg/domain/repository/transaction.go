package repository

import (
	"context"
)

type Repository interface {
	DoInTx(ctx context.Context, f func(context.Context) (interface{}, error)) (interface{}, error)
}
