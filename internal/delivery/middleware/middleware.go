package middleware

import (
	"github.com/rs/zerolog"

	"github.com/willbrid/api-gateway-sql/config"
	"github.com/willbrid/api-gateway-sql/internal/delivery/httpresponse"

	"encoding/base64"
	"net/http"
	"strings"
)

type IAuthMiddleware interface {
	Authenticate(next http.Handler, config *config.Config) http.Handler
}

type AuthMiddleware struct {
	logger zerolog.Logger
}

func NewAuthMiddleware(logger zerolog.Logger) *AuthMiddleware {
	return &AuthMiddleware{logger}
}

func (a *AuthMiddleware) Authenticate(next http.Handler, config *config.Config) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")

		if config.ApiGatewaySQL.Enabled && !strings.HasPrefix(req.RequestURI, "/swagger/") && !strings.HasPrefix(req.RequestURI, "/healthz") {
			if auth == "" {
				a.logger.Error().Msg("no authorization header found")
				_ = httpresponse.SendJSONResponse(resp, http.StatusUnauthorized, "invalid credential", nil)
				return
			}

			if !strings.HasPrefix(auth, "Basic ") {
				a.logger.Error().Msg("invalid authorization header")
				_ = httpresponse.SendJSONResponse(resp, http.StatusUnauthorized, "invalid credential", nil)
				return
			}

			token := strings.TrimPrefix(auth, "Basic ")
			decodedToken, err := base64.StdEncoding.DecodeString(token)
			if err != nil {
				a.logger.Error().Err(err).Msg("failed to decode base64 token")
				_ = httpresponse.SendJSONResponse(resp, http.StatusUnauthorized, "invalid credential", nil)
				return
			}

			credentialParts := strings.SplitN(string(decodedToken), ":", 2)
			username := credentialParts[0]
			password := credentialParts[1]
			if username != config.ApiGatewaySQL.Username || password != config.ApiGatewaySQL.Password {
				a.logger.Error().Msg("invalid username or password")
				_ = httpresponse.SendJSONResponse(resp, http.StatusUnauthorized, "invalid credential", nil)
				return
			}
		}

		next.ServeHTTP(resp, req)
	})
}
