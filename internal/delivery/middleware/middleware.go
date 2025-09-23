package middleware

import (
	"api-gateway-sql/config"
	"api-gateway-sql/pkg/logger"

	"encoding/base64"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler, config *config.Config) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		var auth string = req.Header.Get("Authorization")

		if config.ApiGatewaySQL.Auth.Enabled && !strings.HasPrefix(req.RequestURI, "/swagger/") && !strings.HasPrefix(req.RequestURI, "/healthz") {
			if auth == "" {
				logger.Error("no authorization header found")
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(auth, "Basic ") {
				logger.Error("invalid authorization header")
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(auth, "Basic ")
			decodedToken, err := base64.StdEncoding.DecodeString(token)
			if err != nil {
				logger.Error("failed to decode base64 token - %v", err)
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}

			credentialParts := strings.SplitN(string(decodedToken), ":", 2)
			username := credentialParts[0]
			password := credentialParts[1]
			if username != config.ApiGatewaySQL.Auth.Username || password != config.ApiGatewaySQL.Auth.Password {
				logger.Error("invalid username or password")
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(resp, req)
	})
}
