package httphandler

import (
	"api-gateway-sql/internal/domain"
	"api-gateway-sql/pkg/httpresponse"
	"api-gateway-sql/pkg/logger"
	"encoding/json"

	"context"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	successAPIMessage string = "operation completed successfully"
	failedAPIMessage  string = "Operation ended in failure"
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
	var (
		vars       map[string]string = mux.Vars(req)
		targetName string            = vars["target"]
		ctx        context.Context   = req.Context()
	)

	sqlqueryInput := &domain.SQLQueryInput{
		TargetName: targetName,
		PostParams: make(map[string]any, 0),
	}

	sqlqueryOutput, err := h.Usercases.ISQLQueryUsecase.ExecuteSingle(ctx, sqlqueryInput)
	if err != nil {
		logger.LogError("error: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusInternalServerError, failedAPIMessage, nil)
		return
	}

	httpresponse.SendJSONResponse(resp, http.StatusOK, successAPIMessage, sqlqueryOutput)
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
	var (
		vars       map[string]string = mux.Vars(req)
		targetName string            = vars["target"]
		ctx        context.Context   = req.Context()
	)

	var postParams map[string]any
	if err := json.NewDecoder(req.Body).Decode(&postParams); err != nil {
		logger.LogError("error: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusBadRequest, err.Error(), nil)
		return
	}

	sqlqueryInput := &domain.SQLQueryInput{
		TargetName: targetName,
		PostParams: postParams,
	}

	sqlqueryOutput, err := h.Usercases.ISQLQueryUsecase.ExecuteSingle(ctx, sqlqueryInput)
	if err != nil {
		logger.LogError("error: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusInternalServerError, failedAPIMessage, nil)
		return
	}

	httpresponse.SendJSONResponse(resp, http.StatusOK, successAPIMessage, sqlqueryOutput)
}
