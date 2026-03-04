package httphandler

import (
	"github.com/rs/zerolog"
	"github.com/willbrid/api-gateway-sql/config"
	"github.com/willbrid/api-gateway-sql/internal/usecase"
)

type HTTPHandler struct {
	Usercases *usecase.Usecases
	cfg       *config.Config
	logger    zerolog.Logger
}

func NewHTTPHandler(usercases *usecase.Usecases, cfg *config.Config, logger zerolog.Logger) *HTTPHandler {
	return &HTTPHandler{
		Usercases: usercases,
		cfg:       cfg,
		logger:    logger.With().Str("layer", "handler").Logger(),
	}
}
