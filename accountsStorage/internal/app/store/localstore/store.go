package localstore

import (
	"account_storage/internal/app/store"
	"account_storage/pkg/model"

	"github.com/sirupsen/logrus"
)

type Store struct {
	logger            *logrus.Logger
	accountRepository store.AccountRepository
}

func New(logger *logrus.Logger) *Store {
	return &Store{
		logger: logger,
	}
}

func (store Store) Account() store.AccountRepository {
	if store.accountRepository != nil {
		return store.accountRepository
	}

	return &AccountRepository{
		accounts: make(map[string]model.Account),
		logger:   store.logger,
	}
}
