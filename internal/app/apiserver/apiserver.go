package apiserver

import (
	"account_storage/internal/app/store"
	"account_storage/internal/app/store/localstore"
	"account_storage/internal/app/store/sqlstore"
	"account_storage/pkg/model/account"
	"account_storage/pkg/oc"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	kithttp "github.com/go-kit/kit/transport/http"
	_ "github.com/lib/pq"

	"github.com/go-kit/kit/tracing/opencensus"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *gin.Engine
	logger *logrus.Logger
	store  store.Store
	config *Config
	ctx    context.Context
}

func NewServer(logger *logrus.Logger, ctx context.Context, config *Config) (*server, error) {
	var store store.Store

	switch config.DatabaseType {
	case "sql":
		db, err := newDB(config.DatabaseURL, logger)
		if err != nil {
			return nil, err
		}

		store = sqlstore.New(db, logger)
	case "local":
		store = localstore.New(logger)
	default:
		return nil, fmt.Errorf("unknown database_type %s", config.DatabaseType)
	}

	server := &server{
		router: gin.Default(),
		logger: logger,
		store:  store,
		config: config,
		ctx:    ctx,
	}

	err := server.configureLogger(config.LogLevel)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (server *server) configureLogger(logLevel string) error {

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		server.logger.WithFields(logrus.Fields{
			"package":  "apiserver",
			"function": "configureLogger",
			"error":    err,
			"logLevel": logLevel,
		}).Error("parsing log level failed")

		return err
	}

	server.logger.SetLevel(level)

	return nil
}

func (server *server) Start() error {

	var accountService account.Service
	{
		accountRepository := server.store.Account()
		accountService = account.NewService(accountRepository, server.logger)
	}

	var accountEndpoints account.Endpoints
	{
		accountEndpoints = account.MakeEndpoints(accountService)

		accountEndpoints = account.Endpoints{
			Create:  oc.ServerEndpoint("Create")(accountEndpoints.Create),
			GetByID: oc.ServerEndpoint("GetByID")(accountEndpoints.GetByID),
			Update:  oc.ServerEndpoint("Update")(accountEndpoints.Update),
			Delete:  oc.ServerEndpoint("Delete")(accountEndpoints.Delete),
			GetAll:  oc.ServerEndpoint("GetAll")(accountEndpoints.GetAll),
		}
	}
	var httpHandler http.Handler
	{
		ocTracing := opencensus.HTTPServerTrace()
		serverOptions := []kithttp.ServerOption{ocTracing}
		httpHandler = account.NewGinService(accountEndpoints, serverOptions, server.logger)
	}

	httpServer := &http.Server{
		Addr:    server.config.BindAddres,
		Handler: httpHandler,
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			server.logger.WithFields(logrus.Fields{
				"package":    "apiserver",
				"function":   "Start",
				"error":      err,
				"httpServer": httpServer,
			}).Error("server listen and serve failed")

		}
	}()

	server.logger.Debug("Server started")

	<-server.ctx.Done()
	server.logger.Debug("Starting graceful shotdown")

	shutdownContext, cancel := context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(server.config.ShutdownTimeout))
	defer cancel()

	err := httpServer.Shutdown(shutdownContext)
	if err != nil {
		server.logger.WithFields(logrus.Fields{
			"package":    "apiserver",
			"function":   "Start",
			"error":      err,
			"httpServer": httpServer,
		}).Error("server shutdown failed")

		return err
	}

	server.logger.Debug("Server stopped gracefully")

	return nil
}

func newDB(databaseURL string, logger *logrus.Logger) (*sql.DB, error) {

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"package":  "apiserver",
			"function": "newDB",
			"error":    err,
		}).Error("openning database connection failed")

		return nil, err
	}
	defer func() {
		err := db.Close()
		if err != nil {
			logger.WithFields(logrus.Fields{
				"package":  "apiserver",
				"function": "Start",
				"error":    err,
			}).Error("closing database failed")
		}
	}()

	err = db.Ping()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"package":     "apiserver",
			"function":    "newDB",
			"databaseURL": databaseURL,
			"error":       err,
		}).Error("ping database failed")

		return nil, err
	}
	return db, nil
}
