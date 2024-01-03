package server

import (
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"

	"btc-sbt/logger"
	"btc-sbt/server/middleware"
)

// APIService defines the api service
type APIService struct {
	APIBackend

	Router *gin.Engine

	Logger *logrus.Logger

	Mode string
}

// NewAPIService creates a new APIService instance
func NewAPIService(
	apiBackend APIBackend,
	logger *logrus.Logger,
) *APIService {
	srv := APIService{
		APIBackend: apiBackend,
		Logger:     logger,
	}

	srv.createRouter()

	return &srv
}

// Start starts the api service
func (srv *APIService) Start(listenerAddr string) error {
	srv.Logger.Infof("starting the API service")

	return srv.Router.Run(listenerAddr)
}

// Stop stops the api service
func (srv *APIService) Stop() error {
	srv.Logger.Infof("API service stopped")

	return nil
}

// createRouter creates the router
func (srv *APIService) createRouter() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(middleware.WithBody(), middleware.BodyLogger(logger.Logger))

	r.GET("/api/collections", srv.GetAllSBTs)

	r.GET("/api/collections/:symbol", srv.GetSBTs)
	r.GET("api/sbts", srv.GetSBT)

	r.GET("api/sbts/address/:address", srv.GetOwnedSBTsWrapper)

	r.GET("api/status", srv.Status)

	srv.Router = r
}
