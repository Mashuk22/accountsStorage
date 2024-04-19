package account

import (
	"account_storage/internal/app/store"
	"account_storage/pkg/model"
	"context"
	"log"

	"github.com/sirupsen/logrus"
)

type Service interface {
	Create(ctx context.Context, account model.AccountCreate) (string, error)
	GetByID(ctx context.Context, id string) (model.Account, error)
	Update(ctx context.Context, account model.Account) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]model.Account, error)
	Nginx(ctx context.Context) (string, error)
}

type service struct {
	repository store.AccountRepository
	logger     *logrus.Logger
}

func NewService(repository store.AccountRepository, logger *logrus.Logger) Service {
	return &service{
		repository: repository,
		logger:     logger,
	}
}

// @Summary Create a new account
// @Description Create a new account with the specified details
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body CreateRequest true "Account to create"
// @Success 200 {object} CreateResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /accounts [post]
func (s *service) Create(ctx context.Context, account model.AccountCreate) (string, error) {
	id, err := s.repository.Create(ctx, account)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"package":  "account",
			"function": "Create",
			"error":    err,
			"account":  account,
		}).Error("creating account failed")

		return id, err
	}
	return id, nil
}

// @Summary Get all accounts
// @Description Retrieve a list of all accounts
// @Tags accounts
// @Accept json
// @Produce json
// @Success 200 {object} GetAllResponse "List of accounts"
// @Failure 500 {string} string "Internal Server Error"
// @Router /accounts [get]
func (s *service) GetAll(ctx context.Context) ([]model.Account, error) {
	accounts, err := s.repository.GetAll(ctx)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"package":  "account",
			"function": "GetAll",
			"error":    err,
		}).Error("getting all accounts failed")

		return nil, err
	}
	return accounts, nil
}

// @Summary Get account by ID
// @Description Retrieve an account by its unique identifier
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} GetByIDResponse "Account data"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /accounts/{id} [get]
func (s *service) GetByID(ctx context.Context, id string) (model.Account, error) {
	account, err := s.repository.GetByID(ctx, id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"package":  "account",
			"function": "GetByID",
			"error":    err,
			"id":       id,
		}).Error("getting account by id failed")

		return model.Account{}, err
	}
	return account, nil
}

// @Summary Update an account
// @Description Update an account with specified details
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param account body UpdateRequest true "Account to update"
// @Success 200 {object} UpdateResponse "Updated account data"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /accounts/{id} [put]
func (s *service) Update(ctx context.Context, account model.Account) error {
	err := s.repository.Update(ctx, account)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"package":  "account",
			"function": "Update",
			"error":    err,
			"account":  account,
		}).Error("updating account failed")

		return err
	}
	return nil
}

// @Summary Delete an account
// @Description Delete an account
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID to delete"
// @Success 200 {object} DeleteResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /accounts/{id} [delete]
func (s *service) Delete(ctx context.Context, id string) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"package":  "account",
			"function": "Delete",
			"error":    err,
			"id":       id,
		}).Error("deleting account failed")

		return err
	}
	return nil
}

// @Summary Get data from Nginx
// @Description Makes an HTTP GET request to Nginx and returns the response body as a string.
// @Tags nginx
// @Accept json
// @Produce json
// @Success 200 {string} body "Response body from Nginx"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /nginx [get]
func (s *service) Nginx(ctx context.Context) (string, error) {
	log.Print("start Nginx func in service")
	res, err := s.repository.Nginx(ctx)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"package":  "account",
			"function": "Nginx",
			"error":    err,
		}).Error("Nginx request failed")

		return "", err
	}

	return res, nil
}
