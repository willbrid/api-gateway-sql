package delivery

import (
	"github.com/willbrid/api-gateway-sql/config"
	"github.com/willbrid/api-gateway-sql/internal/delivery/httphandler"
	"github.com/willbrid/api-gateway-sql/internal/delivery/middleware"
	"github.com/willbrid/api-gateway-sql/internal/usecase"
	"github.com/willbrid/api-gateway-sql/pkg/logger"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Handler struct {
	Usecases *usecase.Usecases
	iLogger  logger.ILogger
}

func NewHandler(usecases *usecase.Usecases, iLogger logger.ILogger) *Handler {
	return &Handler{usecases, iLogger}
}

func (h *Handler) InitRouter(router *mux.Router, cfg *config.Config, cfgflag *config.ConfigFlag) {
	router.Use(func(subH http.Handler) http.Handler {
		authMiddleware := middleware.NewAuthMiddleware(h.iLogger)
		return authMiddleware.Authenticate(subH, cfg)
	})

	if cfg.ApiGatewaySQL.EnableSwagger {
		scheme := map[bool]string{true: "https", false: "http"}[cfgflag.EnableHttps]
		swaggerUrl := fmt.Sprintf("%s://localhost:%d/swagger/doc.json", scheme, cfgflag.ListenPort)

		router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
			httpSwagger.URL(swaggerUrl),
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
			httpSwagger.DomID("swagger-ui"),
		)).Methods("GET")
	}

	httphandler := httphandler.NewHTTPHandler(h.Usecases, cfg, h.iLogger)

	router.HandleFunc("/healthz", httphandler.HandleHealthCheck).Methods("GET")
	router.HandleFunc("/api-gateway-sql/blocks/{uid}", httphandler.ApiGetBlockHandler).Methods("GET")
	router.HandleFunc("/api-gateway-sql/batchstats", httphandler.ApiListBatchStatsHandler).Methods("GET")
	router.HandleFunc("/api-gateway-sql/batchstats/{uid}", httphandler.ApiGetBatchStatHandler).Methods("GET")
	router.HandleFunc("/api-gateway-sql/batchstats/{uid}/blocks", httphandler.ApiListBlocksByBatchStatHandler).Methods("GET")
	router.HandleFunc("/api-gateway-sql/batchstats/{uid}/completed", httphandler.ApiMarkCompletedBatchStatHandler).Methods("GET")
	router.HandleFunc("/api-gateway-sql/{target}", httphandler.ApiGetSqlHandler).Methods("GET")
	router.HandleFunc("/api-gateway-sql/{target}", httphandler.ApiPostSqlHandler).Methods("POST")
	router.HandleFunc("/api-gateway-sql/{target}/batch", httphandler.ApiPostSqlBatchHandler).Methods("POST")
	router.HandleFunc("/api-gateway-sql/{datasource}/init", httphandler.ApiPostInitDatabase).Methods("POST")
}
