package account

import (
	"account_storage/pkg/model"
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	Create  endpoint.Endpoint
	GetByID endpoint.Endpoint
	Update  endpoint.Endpoint
	Delete  endpoint.Endpoint
	GetAll  endpoint.Endpoint
	Nginx   endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create:  makeCreateEndpoint(s),
		GetByID: makeGetByIDEndpoint(s),
		Update:  makeUpdateEndpoint(s),
		Delete:  makeDeleteEndpoint(s),
		GetAll:  makeGetAllEndpoint(s),
		Nginx:   makeNginxEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		id, err := s.Create(ctx, req.Account)
		return CreateResponse{ID: id, Err: err}, nil
	}
}

func makeGetAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		accounts, err := s.GetAll(ctx)
		return GetAllResponse{Accounts: accounts, Err: err}, nil
	}
}

func makeGetByIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		accountRes, err := s.GetByID(ctx, req.ID)
		return GetByIDResponse{Account: accountRes, Err: err}, nil
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.Account)
		err := s.Update(ctx, req)
		return UpdateResponse{Err: err}, nil
	}
}
func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)
		err := s.Delete(ctx, req.ID)
		return DeleteResponse{Err: err}, nil
	}
}

func makeNginxEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.Nginx(ctx)
		return res, err
	}
}

type CreateRequest struct {
	Account model.AccountCreate `json:"account"`
}

type CreateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

type GetByIDRequest struct {
	ID string `json:"id"`
}

type GetByIDResponse struct {
	Account model.Account `json:"account"`
	Err     error         `json:"error,omitempty"`
}

type GetAllRequest struct {
}

type GetAllResponse struct {
	Accounts []model.Account `json:"accounts"`
	Err      error           `json:"error,omitempty"`
}

type UpdateRequest struct {
	Account model.AccountUpdate `json:"account"`
}

type UpdateResponse struct {
	Err error `json:"error,omitempty"`
}

type DeleteRequest struct {
	ID string `json:"id"`
}

type DeleteResponse struct {
	Err error `json:"error,omitempty"`
}
