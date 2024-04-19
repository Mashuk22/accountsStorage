package sqlstore

import (
	"account_storage/internal/app/store"
	"database/sql"

	"github.com/sirupsen/logrus"
)

type Store struct {
	db                *sql.DB
	logger            *logrus.Logger
	accountRepository store.AccountRepository
}

func New(db *sql.DB, logger *logrus.Logger) *Store {
	return &Store{
		db:     db,
		logger: logger,
	}
}

func (store Store) Account() store.AccountRepository {
	if store.accountRepository != nil {
		return store.accountRepository
	}

	return &AccountRepository{
		db:     store.db,
		logger: store.logger,
	}
}
