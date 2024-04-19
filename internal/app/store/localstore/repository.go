package localstore

import (
	"account_storage/pkg/model"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AccountRepository struct {
	sync.Mutex
	accounts map[string]model.Account
	logger   *logrus.Logger
}

func (accountRepository *AccountRepository) Create(ctx context.Context, accountCreate model.AccountCreate) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	accountRepository.Lock()
	defer accountRepository.Unlock()

	accountID := uuid.New()
	accountCreatedAt := time.Now()

	account := model.Account{
		ID:                    accountID,
		Name:                  accountCreate.Name,
		AccountType:           accountCreate.AccountType,
		Login:                 accountCreate.Login,
		Password:              accountCreate.Password,
		Email:                 accountCreate.Email,
		EmailPassword:         accountCreate.EmailPassword,
		RecoveryEmail:         accountCreate.RecoveryEmail,
		RecoveryEmailPassword: accountCreate.RecoveryEmailPassword,
		Cookie:                accountCreate.Cookie,
		Status:                accountCreate.Status,
		CreatedAt:             accountCreatedAt,
	}

	stringAccountID := accountID.String()
	accountRepository.accounts[stringAccountID] = account

	return stringAccountID, nil
}

func (accountRepository *AccountRepository) GetByID(ctx context.Context, id string) (model.Account, error) {
	select {
	case <-ctx.Done():
		return model.Account{}, ctx.Err()
	default:
	}

	accountRepository.Lock()
	defer accountRepository.Unlock()

	account, ok := accountRepository.accounts[id]
	if !ok {
		return model.Account{}, fmt.Errorf("no account with id %s", id)
	}

	return account, nil
}

func (accountRepository *AccountRepository) Update(ctx context.Context, accountUpdate model.Account) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	accountRepository.Lock()
	defer accountRepository.Unlock()

	strID := accountUpdate.ID.String()
	account, ok := accountRepository.accounts[strID]
	if !ok {
		return fmt.Errorf("no account with id %s", strID)
	}

	if accountUpdate.Name != "" {
		account.Name = accountUpdate.Name
	}
	if accountUpdate.AccountType != "" {
		account.AccountType = accountUpdate.AccountType
	}
	if accountUpdate.Login != "" {
		account.Login = accountUpdate.Login
	}
	if accountUpdate.Password != "" {
		account.Password = accountUpdate.Password
	}
	if accountUpdate.Email != "" {
		account.Email = accountUpdate.Email
	}
	if accountUpdate.EmailPassword != "" {
		account.EmailPassword = accountUpdate.EmailPassword
	}
	if accountUpdate.RecoveryEmail != "" {
		account.RecoveryEmail = accountUpdate.RecoveryEmail
	}
	if accountUpdate.RecoveryEmailPassword != "" {
		account.RecoveryEmailPassword = accountUpdate.RecoveryEmailPassword
	}
	if accountUpdate.Cookie != "" {
		account.Cookie = accountUpdate.Cookie
	}
	if accountUpdate.Status != "" {
		account.Status = accountUpdate.Status
	}

	accountRepository.accounts[strID] = account

	return nil
}

func (accountRepository *AccountRepository) Delete(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	accountRepository.Lock()
	defer accountRepository.Unlock()

	_, ok := accountRepository.accounts[id]
	if !ok {
		return fmt.Errorf("no account with id %s", id)
	}

	delete(accountRepository.accounts, id)

	return nil
}

func (accountRepository *AccountRepository) GetAll(ctx context.Context) ([]model.Account, error) {
	select {
	case <-ctx.Done():
		return []model.Account{}, ctx.Err()
	default:
	}

	accountRepository.Lock()
	defer accountRepository.Unlock()

	accounts := make([]model.Account, 0, len(accountRepository.accounts))
	for _, account := range accountRepository.accounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (accountRepository *AccountRepository) Nginx(ctx context.Context) (string, error) {
	log.Print("start Nginx func in repository")

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:80", nil)
	if err != nil {
		return "", fmt.Errorf("error http.NewRequestWithContext: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error client.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get a successful response: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}
