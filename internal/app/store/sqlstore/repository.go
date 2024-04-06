package sqlstore

import (
	"account_storage/pkg/model"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AccountRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func (accountRepository *AccountRepository) Create(ctx context.Context, accountCreate model.AccountCreate) (string, error) {
	query := `INSERT INTO accounts (id, name, account_type, login, password, email, email_password, recovery_email, recovery_email_password, cookie, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`

	accountID := uuid.New()
	accountCreatedAt := time.Now()

	var id string
	err := accountRepository.db.QueryRowContext(ctx, query,
		accountID,
		accountCreate.Name,
		accountCreate.AccountType,
		accountCreate.Login,
		accountCreate.Password,
		accountCreate.Email,
		accountCreate.EmailPassword,
		accountCreate.RecoveryEmail,
		accountCreate.RecoveryEmailPassword,
		accountCreate.Cookie,
		accountCreate.Status,
		accountCreatedAt).Scan(&id)

	if err != nil {
		accountRepository.logger.WithError(err).Error("Failed to create account")
		return "", fmt.Errorf("error creating account: %w", err)
	}

	return id, nil
}

func (accountRepository *AccountRepository) GetByID(ctx context.Context, id string) (model.Account, error) {
	query := `SELECT id, name, account_type, login, password, email, email_password, recovery_email, recovery_email_password, cookie, status, created_at
		FROM accounts WHERE id = $1`

	var account model.Account
	err := accountRepository.db.QueryRowContext(ctx, query, id).Scan(
		&account.ID,
		&account.Name,
		&account.AccountType,
		&account.Login,
		&account.Password,
		&account.Email,
		&account.EmailPassword,
		&account.RecoveryEmail,
		&account.RecoveryEmailPassword,
		&account.Cookie,
		&account.Status,
		&account.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			accountRepository.logger.WithError(err).Error("Failed to get account by id")
			return model.Account{}, fmt.Errorf("no account with id %s", id)
		}
		accountRepository.logger.WithError(err).Error("Failed to get account by id")
		return model.Account{}, fmt.Errorf("error getting account by id: %w", err)
	}

	return account, nil
}

func (accountRepository *AccountRepository) Update(ctx context.Context, account model.Account) error {
	query := `UPDATE accounts SET name = $2, account_type = $3, login = $4, password = $5, email = $6, email_password = $7, 
		recovery_email = $8, recovery_email_password = $9, cookie = $10, status = $11
		WHERE id = $1`

	_, err := accountRepository.db.ExecContext(ctx, query,
		account.ID,
		account.Name,
		account.AccountType,
		account.Login,
		account.Password,
		account.Email,
		account.EmailPassword,
		account.RecoveryEmail,
		account.RecoveryEmailPassword,
		account.Cookie,
		account.Status,
	)

	accountRepository.logger.WithError(err).Error("Failed to update account")
	return fmt.Errorf("error updating account with id %s: %w", account.ID, err)
}

func (accountRepository *AccountRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM accounts WHERE id = $1`

	_, err := accountRepository.db.ExecContext(ctx, query, id)

	if err != nil {
		accountRepository.logger.WithError(err).Error("Failed to delete account")
		return fmt.Errorf("error deleting account with id %s: %w", id, err)
	}

	return nil
}

func (accountRepository *AccountRepository) GetAll(ctx context.Context) ([]model.Account, error) {
	query := "SELECT id, name, account_type, login, password, email, email_password, recovery_email, recovery_email_password, cookie, status, created_at FROM accounts"

	rows, err := accountRepository.db.QueryContext(ctx, query)
	if err != nil {
		accountRepository.logger.WithError(err).Error("Failed to get all accounts")
		return nil, fmt.Errorf("error getting all accounts: %w", err)
	}
	defer rows.Close()

	var accounts []model.Account

	for rows.Next() {
		var account model.Account
		err := rows.Scan(
			&account.ID,
			&account.Name,
			&account.AccountType,
			&account.Login,
			&account.Password,
			&account.Email,
			&account.EmailPassword,
			&account.RecoveryEmail,
			&account.RecoveryEmailPassword,
			&account.Cookie,
			&account.Status,
			&account.CreatedAt,
		)
		if err != nil {
			accountRepository.logger.WithError(err).Error("Failed to get all accounts")
			return nil, fmt.Errorf("error getting all accounts: %w", err)
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		accountRepository.logger.WithError(err).Error("Failed to get all accounts")
		return nil, fmt.Errorf("error getting all accounts: %w", err)
	}

	return accounts, nil
}
