package middleware

import (
	"github.com/willbrid/api-gateway-sql/config"
	"github.com/willbrid/api-gateway-sql/internal/delivery/httpresponse"
	"github.com/willbrid/api-gateway-sql/pkg/logger"

	"encoding/base64"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	iLogger logger.ILogger
}

func NewAuthMiddleware(iLogger logger.ILogger) *AuthMiddleware {
	return &AuthMiddleware{iLogger}
}

func (a *AuthMiddleware) Authenticate(next http.Handler, config *config.Config) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		var auth string = req.Header.Get("Authorization")

		if config.ApiGatewaySQL.Auth.Enabled && !strings.HasPrefix(req.RequestURI, "/swagger/") && !strings.HasPrefix(req.RequestURI, "/healthz") {
			if auth == "" {
				a.iLogger.Error("no authorization header found")
				httpresponse.SendJSONResponse(resp, http.StatusUnauthorized, "invalid credential", nil)
				return
			}

			if !strings.HasPrefix(auth, "Basic ") {
				a.iLogger.Error("invalid authorization header")
				httpresponse.SendJSONResponse(resp, http.StatusUnauthorized, "invalid credential", nil)
				return
			}

			token := strings.TrimPrefix(auth, "Basic ")
			decodedToken, err := base64.StdEncoding.DecodeString(token)
			if err != nil {
				a.iLogger.Error("failed to decode base64 token - %v", err)
				httpresponse.SendJSONResponse(resp, http.StatusUnauthorized, "invalid credential", nil)
				return
			}

			credentialParts := strings.SplitN(string(decodedToken), ":", 2)
			username := credentialParts[0]
			password := credentialParts[1]
			if username != config.ApiGatewaySQL.Auth.Username || password != config.ApiGatewaySQL.Auth.Password {
				a.iLogger.Error("invalid username or password")
				httpresponse.SendJSONResponse(resp, http.StatusUnauthorized, "invalid credential", nil)
				return
			}
		}

		next.ServeHTTP(resp, req)
	})
}
