package httphandler

import (
	"api-gateway-sql/internal/domain"
	"api-gateway-sql/pkg/httpresponse"
	"api-gateway-sql/pkg/logger"

	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	successAPIMessage              string = "operation completed successfully"
	failedAPIMessage               string = "operation ended in failure"
	errUnableToReadSQLFile         string = "unable to read the sql file content"
	errUnableToReadCSVFile         string = "unable to read the csv file content"
	errUnableToExecuteInitSqlQuery string = "unable to execute the init sql query"
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
		logger.Error("error: %s", err.Error())
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
		logger.Error("error: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusBadRequest, err.Error(), nil)
		return
	}

	sqlqueryInput := &domain.SQLQueryInput{
		TargetName: targetName,
		PostParams: postParams,
	}

	sqlqueryOutput, err := h.Usercases.ISQLQueryUsecase.ExecuteSingle(ctx, sqlqueryInput)
	if err != nil {
		logger.Error("error: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusInternalServerError, failedAPIMessage, nil)
		return
	}

	httpresponse.SendJSONResponse(resp, http.StatusOK, successAPIMessage, sqlqueryOutput)
}

// ApiPostInitDatabase godoc
// @Summary      Initialize Database
// @Description  Initialize Database by providing a sql query file
// @Tags         apisql
// @Accept       json
// @Produce      json
// @Param        datasource  path  string  true  "Datasource Name"
// @Param        sqlfile  formData  file  true  "SQL Data to upload"
// @Success      200  {object}  httpresponse.HTTPResp
// @Failure      400  {object}  httpresponse.HTTPResp
// @Failure      500  {object}  httpresponse.HTTPResp
// @Security     BasicAuth
// @Router       /api-gateway-sql/{datasource}/init [post]
func (h *HTTPHandler) ApiPostInitDatabase(resp http.ResponseWriter, req *http.Request) {
	var (
		vars           map[string]string = mux.Vars(req)
		datasourceName string            = vars["datasource"]
		ctx            context.Context   = req.Context()
	)

	file, _, err := req.FormFile("sqlfile")
	if err != nil {
		logger.Error("error: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusBadRequest, errUnableToReadSQLFile, nil)
		return
	}

	sqlBytes, err := io.ReadAll(file)
	if err != nil {
		logger.Error("error: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusBadRequest, errUnableToReadSQLFile, nil)
		return
	}

	sqlInitDatabaseInput := &domain.SQLInitDatabaseInput{
		Datasource:     datasourceName,
		SQLFileContent: string(sqlBytes),
	}

	if err := h.Usercases.ISQLInitDatabaseUsecase.ExecuteInit(ctx, sqlInitDatabaseInput); err != nil {
		logger.Error("error: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusInternalServerError, errUnableToExecuteInitSqlQuery, nil)
		return
	}

	httpresponse.SendJSONResponse(resp, http.StatusOK, successAPIMessage, nil)
}

// ApiPostSqlBatchHandler godoc
// @Summary      Execute batch sql query
// @Description  Execute batch sql query with values from a csv file
// @Tags         apisql
// @Accept       json
// @Produce      json
// @Param        target path  string  true  "Target Name"
// @Param        csvfile  formData  file  true  "CSV Data to import"
// @Success      200  {object}  httpresponse.HTTPResp
// @Failure      400  {object}  httpresponse.HTTPResp
// @Failure      500  {object}  httpresponse.HTTPResp
// @Security     BasicAuth
// @Router       /api-gateway-sql/{target}/batch [post]
func (h *HTTPHandler) ApiPostSqlBatchHandler(resp http.ResponseWriter, req *http.Request) {
	var (
		vars       map[string]string = mux.Vars(req)
		targetName string            = vars["target"]
		ctx        context.Context   = req.Context()
	)

	csvfile, _, err := req.FormFile("csvfile")
	if err != nil {
		logger.Error("error: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusBadRequest, errUnableToReadCSVFile, nil)
		return
	}

	sqlBatchQueryInput := &domain.SQLBatchQueryInput{
		TargetName: targetName,
		File:       csvfile,
	}

	go func() {
		if err := h.Usercases.ISQLBatchQueryUsecase.ExecuteBatch(ctx, sqlBatchQueryInput); err != nil {
			logger.Error("error: %s", err.Error())
		}
	}()

	httpresponse.SendJSONResponse(resp, http.StatusOK, httpresponse.HTTPStatusOKMessage, nil)
}
