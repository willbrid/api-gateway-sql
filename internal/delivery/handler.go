package delivery

import (
	"api-gateway-sql/config"
	"api-gateway-sql/internal/delivery/httphandler"
	"api-gateway-sql/internal/delivery/middleware"
	"api-gateway-sql/internal/usecase"

	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	Usercases *usecase.Usecases
}

func NewHandler(usercases *usecase.Usecases) *Handler {
	return &Handler{usercases}
}

func (h *Handler) InitRouter(router *mux.Router, cfg *config.Config) {
	router.Use(func(subH http.Handler) http.Handler {
		return middleware.AuthMiddleware(subH, cfg)
	})

	httphandler := httphandler.NewHTTPHandler(h.Usercases, cfg)

	router.HandleFunc("/healthz", httphandler.HandleHealthCheck).Methods("GET")
	router.HandleFunc("/api-gateway-sql/{target}", httphandler.ApiGetSqlHandler).Methods("GET")
	router.HandleFunc("/api-gateway-sql/{target}", httphandler.ApiPostSqlHandler).Methods("POST")
}
