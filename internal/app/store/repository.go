package store

import (
	"account_storage/pkg/model"
	"context"
)

type AccountRepository interface {
	Create(ctx context.Context, account model.AccountCreate) (string, error)
	GetByID(ctx context.Context, id string) (model.Account, error)
	Update(ctx context.Context, account model.Account) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]model.Account, error)
}
