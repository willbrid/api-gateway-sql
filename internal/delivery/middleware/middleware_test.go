package middleware_test

import (
	"github.com/willbrid/api-gateway-sql/config"
	"github.com/willbrid/api-gateway-sql/internal/delivery/middleware"

	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var yamlConfig []byte = []byte(`
---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "sqlite"
    dbname: "/tmp/xxxxx"
    timeout: "10s"
  targets:
  - name: xxxxx
    data_source_name: xxxxx
    datafields: ""
    sql: "select * from students"
`)

func triggerTest(t *testing.T, statusCode int, credential string) {
	v := viper.New()
	v.SetConfigType("yaml")

	if err := v.ReadConfig(bytes.NewBuffer([]byte(yamlConfig))); err != nil {
		t.Fatal(err.Error())
	}

	var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())
	configLoaded, err := config.LoadConfig(v, validate)
	if err != nil {
		t.Fatal(err.Error())
	}

	req, err := http.NewRequest("GET", "/api-gateway-sql/xxxxx", nil)
	if err != nil {
		t.Fatal(err)
	}

	if credential != "" {
		req.Header.Add("Authorization", credential)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api-gateway-sql/{targetname}", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
	}).Methods("GET")
	router.Use(func(next http.Handler) http.Handler {
		return middleware.AuthMiddleware(next, configLoaded)
	})
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != statusCode {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, statusCode)
	}
}

func TestAuthentication_NoAuthorizationHeader(t *testing.T) {
	t.Parallel()

	triggerTest(t, http.StatusUnauthorized, "")
}

func TestAuthentication_InvalidAuthorizationHeader(t *testing.T) {
	t.Parallel()

	triggerTest(t, http.StatusUnauthorized, "xxxxx")
}

func TestAuthentication_FailedToDecodeBase64Token(t *testing.T) {
	t.Parallel()

	triggerTest(t, http.StatusUnauthorized, "Basic xxxxx")
}

func TestAuthentication_InvalidUsernameOrPassword(t *testing.T) {
	t.Parallel()

	triggerTest(t, http.StatusUnauthorized, "Basic eHh4eHg6eHh4")
}

func TestAuthentication_CorrectUsernameOrPassword(t *testing.T) {
	t.Parallel()

	triggerTest(t, http.StatusOK, "Basic eHh4eHg6eHh4eHh4eHg=")
}
