package account

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	_ "account_storage/docs"
	"account_storage/pkg/logadapter"
	"account_storage/pkg/model"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	ErrBadRouting = errors.New("bad routing")
)

type GinContextKey struct{}

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), GinContextKey{}, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func NewGinService(
	svcEndpoints Endpoints, options []kithttp.ServerOption, logger *logrus.Logger) *gin.Engine {
	router := gin.Default()
	router.Use(GinContextToContextMiddleware())
	logrusAdapter := logadapter.NewLogrusAdapter(logger)
	errorLogger := kithttp.ServerErrorLogger(logrusAdapter)
	errorEncoder := kithttp.ServerErrorEncoder(encodeErrorResponse)
	options = append(options, errorLogger, errorEncoder)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/nginx", nginxHandler(logger))

	router.POST("/accounts", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		kithttp.NewServer(
			svcEndpoints.Create,
			decodeCreateRequest(logger),
			encodeResponse(logger),
			options...,
		).ServeHTTP(w, r)
	}))

	router.GET("/accounts", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		kithttp.NewServer(
			svcEndpoints.GetAll,
			decodeGetAllRequest,
			encodeResponse(logger),
			options...,
		).ServeHTTP(w, r)
	}))

	router.GET("/accounts/:id", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		kithttp.NewServer(
			svcEndpoints.GetByID,
			decodeGetByIDRequest,
			encodeResponse(logger),
			options...,
		).ServeHTTP(w, r)
	}))

	router.PUT("/accounts/:id", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		kithttp.NewServer(
			svcEndpoints.Update,
			decodeUpdateRequest(logger),
			encodeResponse(logger),
			options...,
		).ServeHTTP(w, r)
	}))

	router.DELETE("/accounts/:id", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		kithttp.NewServer(
			svcEndpoints.Delete,
			decodeDeleteRequest,
			encodeResponse(logger),
			options...,
		).ServeHTTP(w, r)
	}))

	return router
}

func nginxHandler(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := http.Get("http://accounts-storage-nginx")
		if err != nil {
			logger.Errorf("Failed to get response from nginx: %v", err)
			c.String(http.StatusInternalServerError, "Error contacting Nginx server")
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Errorf("Failed to read response from nginx: %v", err)
			c.String(http.StatusInternalServerError, "Error reading response from Nginx server")
			return
		}

		c.String(resp.StatusCode, string(body))
	}
}

func decodeCreateRequest(logger *logrus.Logger) kithttp.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		var req CreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.WithFields(logrus.Fields{
				"package":  "account",
				"function": "decodeCreateRequest",
				"error":    err,
			}).Error("decoding from json failed")

			return nil, err
		}

		return req, nil
	}

}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	ginCtx, ok := r.Context().Value(GinContextKey{}).(*gin.Context)
	if !ok {
		return nil, errors.New("could not retrieve gin.Context")
	}

	id := ginCtx.Param("id")
	if id == "" {
		return nil, ErrBadRouting
	}
	return GetByIDRequest{ID: id}, nil
}

func decodeGetAllRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return GetAllRequest{}, nil
}

func decodeUpdateRequest(logger *logrus.Logger) kithttp.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		ginCtx, ok := r.Context().Value(GinContextKey{}).(*gin.Context)
		if !ok {
			logger.WithFields(logrus.Fields{
				"package":  "account",
				"function": "decodeUpdateRequest",
				"error":    "could not retrieve gin.Context",
			}).Error("server listen and serve failed")

			return nil, errors.New("could not retrieve gin.Context")
		}

		idStr := ginCtx.Param("id")
		if idStr == "" {
			return nil, ErrBadRouting
		}

		idUUID, err := uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}

		var req UpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.WithFields(logrus.Fields{
				"package":  "account",
				"function": "decodeUpdateRequest",
				"error":    err,
			}).Error("decoding from json failed")
			return nil, err
		}

		account := model.Account{
			ID:                    idUUID,
			Name:                  req.Account.Name,
			AccountType:           req.Account.AccountType,
			Login:                 req.Account.Login,
			Password:              req.Account.Password,
			Email:                 req.Account.Email,
			EmailPassword:         req.Account.EmailPassword,
			RecoveryEmail:         req.Account.RecoveryEmail,
			RecoveryEmailPassword: req.Account.RecoveryEmailPassword,
			Cookie:                req.Account.Cookie,
			Status:                req.Account.Status,
		}

		return account, nil
	}
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req DeleteRequest

	ginCtx, ok := r.Context().Value(GinContextKey{}).(*gin.Context)
	if !ok {
		return nil, errors.New("could not retrieve gin.Context")
	}

	idStr := ginCtx.Param("id")
	if idStr == "" {
		return nil, ErrBadRouting
	}

	req.ID = idStr

	return req, nil
}

func encodeResponse(logger *logrus.Logger) kithttp.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		logger.WithFields(logrus.Fields{
			"package":  "account",
			"function": "encodeResponse",
		}).Debug("start encodeResponse function")

		if e, ok := response.(errorer); ok && e.error() != nil {
			logger.Errorf("Handling error: %v", e.error())
			encodeErrorResponse(ctx, e.error(), w)
			return nil
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Errorf("Error encoding JSON response: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
		logger.Debug("Successfully encoded response to JSON")
		return nil
	}

}

type errorer interface {
	error() error
}

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	default:
		return http.StatusInternalServerError
	}
}
