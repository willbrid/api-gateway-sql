package httphandler

import (
	"net/http"
)

// HandleHealthCheck godoc
// @Summary      Health check for application
// @Description  Trigger SQL query without params
// @Tags         apisql
// @Accept       json
// @Produce      json
// @Success      204  {object}  httpresponse.HTTPResp
// @Router       /healthz [get]
func (h *HTTPHandler) HandleHealthCheck(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusNoContent)
}

// ApiGetSqlHandler godoc
// @Summary      Get SQL Query result
// @Description  Trigger SQL query without params
// @Tags         apisql
// @Accept       json
// @Produce      json
// @Param        target  path  string  true  "Target Name"
// @Success      200  {object}  httpresponse.HTTPResp
// @Failure      400  {object}  httpresponse.HTTPResp
// @Failure      500  {object}  httpresponse.HTTPResp
// @Security     BasicAuth
// @Router       /api-gateway-sql/{target} [get]
func (h *HTTPHandler) ApiGetSqlHandler(resp http.ResponseWriter, req *http.Request) {

}

// ApiPostSqlHandler godoc
// @Summary      Get SQL Query result
// @Description  Trigger SQL query with params
// @Tags         apisql
// @Accept       json
// @Produce      json
// @Param        target path  string  true  "Target Name"
// @Param        data  body  map[string]interface{}  true  "Data to send"
// @Success      200  {object}  httpresponse.HTTPResp
// @Failure      400  {object}  httpresponse.HTTPResp
// @Failure      500  {object}  httpresponse.HTTPResp
// @Security     BasicAuth
// @Router       /api-gateway-sql/{target} [post]
func (h *HTTPHandler) ApiPostSqlHandler(resp http.ResponseWriter, req *http.Request) {

}
