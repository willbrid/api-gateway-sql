package httphandler

import (
	"github.com/willbrid/api-gateway-sql/config"
	"github.com/willbrid/api-gateway-sql/internal/usecase"
)

type HTTPHandler struct {
	Usercases *usecase.Usecases
	cfg       *config.Config
}

func NewHTTPHandler(usercases *usecase.Usecases, cfg *config.Config) *HTTPHandler {
	return &HTTPHandler{usercases, cfg}
}
