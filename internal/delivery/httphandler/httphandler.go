package httphandler

import (
	"github.com/willbrid/api-gateway-sql/config"
	"github.com/willbrid/api-gateway-sql/internal/usecase"
	"github.com/willbrid/api-gateway-sql/pkg/logger"
)

type HTTPHandler struct {
	Usercases *usecase.Usecases
	cfg       *config.Config
	iLogger   logger.ILogger
}

func NewHTTPHandler(usercases *usecase.Usecases, cfg *config.Config, iLogger logger.ILogger) *HTTPHandler {
	return &HTTPHandler{usercases, cfg, iLogger}
}
