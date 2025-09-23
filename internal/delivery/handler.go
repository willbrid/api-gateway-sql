package delivery

import (
	"api-gateway-sql/config"
	"api-gateway-sql/internal/delivery/httphandler"
	"api-gateway-sql/internal/delivery/middleware"
	"api-gateway-sql/internal/usecase"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Handler struct {
	Usecases *usecase.Usecases
}

func NewHandler(usecases *usecase.Usecases) *Handler {
	return &Handler{usecases}
}

func (h *Handler) InitRouter(router *mux.Router, cfg *config.Config, cfgflag *config.ConfigFlag) {
	router.Use(func(subH http.Handler) http.Handler {
		return middleware.AuthMiddleware(subH, cfg)
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

	httphandler := httphandler.NewHTTPHandler(h.Usecases, cfg)

	router.HandleFunc("/healthz", httphandler.HandleHealthCheck).Methods("GET")
	router.HandleFunc("/api-gateway-sql/{target}", httphandler.ApiGetSqlHandler).Methods("GET")
	router.HandleFunc("/api-gateway-sql/{target}", httphandler.ApiPostSqlHandler).Methods("POST")
	router.HandleFunc("/api-gateway-sql/{datasource}/init", httphandler.ApiPostInitDatabase).Methods("POST")
}
