package httphandler

import (
	"github.com/willbrid/api-gateway-sql/internal/dto"
	"github.com/willbrid/api-gateway-sql/internal/dto/paginator"
	"github.com/willbrid/api-gateway-sql/pkg/httpresponse"
	"github.com/willbrid/api-gateway-sql/pkg/logger"

	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	successAPIMessage              string = "operation completed successfully"
	failedAPIMessage               string = "operation ended in failure"
	errUnableToReadSQLFile         string = "unable to read the sql file content"
	errUnableToReadCSVFile         string = "unable to read the csv file content"
	errUnableToExecuteInitSqlQuery string = "unable to execute the init sql query"
	errThereIsAnUnCompletedBatch   string = "there is an uncompleted batch"
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

	sqlqueryInput := &dto.SQLQueryInput{
		TargetName: targetName,
		PostParams: make(map[string]any, 0),
	}

	sqlqueryOutput, err := h.Usercases.ISQLQueryUsecase.ExecuteSingle(ctx, sqlqueryInput)
	if err != nil {
		logger.Error("failed to decode post params: %s", err.Error())
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
		logger.Error("failed to decode post params: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusBadRequest, err.Error(), nil)
		return
	}

	sqlqueryInput := &dto.SQLQueryInput{
		TargetName: targetName,
		PostParams: postParams,
	}

	sqlqueryOutput, err := h.Usercases.ISQLQueryUsecase.ExecuteSingle(ctx, sqlqueryInput)
	if err != nil {
		logger.Error("failed to execute single sql query: %s", err.Error())
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
		logger.Error("failed to read sql file: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusBadRequest, errUnableToReadSQLFile, nil)
		return
	}

	sqlBytes, err := io.ReadAll(file)
	if err != nil {
		logger.Error("failed to read sql file: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusBadRequest, errUnableToReadSQLFile, nil)
		return
	}

	sqlInitDatabaseInput := &dto.SQLInitDatabaseInput{
		Datasource:     datasourceName,
		SQLFileContent: string(sqlBytes),
	}

	if err := h.Usercases.ISQLQueryUsecase.ExecuteInit(ctx, sqlInitDatabaseInput); err != nil {
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
		logger.Error("failed to read csv file: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusBadRequest, errUnableToReadCSVFile, nil)
		return
	}

	if isAllBatchStatClosed, err := h.Usercases.IBatchStatUsecase.IsAllBatchStatClosed(ctx); err != nil {
		logger.Error("an error occurred: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusInternalServerError, httpresponse.HTTPStatusInternalServerErrorMessage, nil)
		return
	} else if isAllBatchStatClosed == false {
		httpresponse.SendJSONResponse(resp, http.StatusBadRequest, errThereIsAnUnCompletedBatch, nil)
		return
	}

	sqlBatchQueryInput := &dto.SQLBatchQueryInput{
		TargetName: targetName,
		File:       csvfile,
	}

	go func() {
		ctx := context.Background()
		if err := h.Usercases.ISQLBatchQueryUsecase.ExecuteBatch(ctx, sqlBatchQueryInput); err != nil {
			logger.Error("failed to process batch: %s", err.Error())
		}
	}()

	httpresponse.SendJSONResponse(resp, http.StatusOK, httpresponse.HTTPStatusOKMessage, nil)
}

// ApiListBatchStatsHandler godoc
// @Summary      Get BatchStats result
// @Description  Get a paginated list of BatchStat
// @Tags         apisql
// @Accept       json
// @Produce      json
// @Param        page_num query  int  false  "Page number" default(1)
// @Param        page_size query int  false  "Page size" default(20)
// @Success      200  {object}  httpresponse.HTTPResp
// @Failure      400  {object}  httpresponse.HTTPResp
// @Failure      500  {object}  httpresponse.HTTPResp
// @Security     BasicAuth
// @Router       /api-gateway-sql/batchstats [get]
func (h *HTTPHandler) ApiListBatchStatsHandler(resp http.ResponseWriter, req *http.Request) {
	var (
		queries  url.Values = req.URL.Query()
		pageNum  int
		pageSize int
		ctx      context.Context = req.Context()
		err      error
	)

	pageNum, err = strconv.Atoi(queries.Get("page_num"))
	if err != nil {
		logger.Error("failed to handle page_num param query: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusBadRequest, "Unable to handle page_num", nil)
		return
	}
	pageSize, err = strconv.Atoi(queries.Get("page_size"))
	if err != nil {
		logger.Error("failed to handle page_size param query: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusBadRequest, "Unable to handle page_size", nil)
		return
	}

	pageRequest := paginator.NewPageRequest(pageNum, pageSize)
	var statsPaginatorResponse *paginator.PageResponse

	statsPaginatorResponse, err = h.Usercases.IBatchStatUsecase.ListBatchStats(ctx, pageRequest)
	if err != nil {
		logger.Error("failed to list batchstat: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusInternalServerError, httpresponse.HTTPStatusInternalServerErrorMessage, nil)
		return
	}

	httpresponse.SendJSONResponse(resp, http.StatusOK, httpresponse.HTTPStatusOKMessage, statsPaginatorResponse)
}

// ApiGetBatchStatHandler godoc
// @Summary      Get BatchStat result
// @Description  Get a BatchStat by uid
// @Tags         apisql
// @Accept       json
// @Produce      json
// @Param        uid path string  true  "uid"
// @Success      200  {object}  httpresponse.HTTPResp
// @Failure      500  {object}  httpresponse.HTTPResp
// @Security     BasicAuth
// @Router       /api-gateway-sql/batchstats/{uid} [get]
func (h *HTTPHandler) ApiGetBatchStatHandler(resp http.ResponseWriter, req *http.Request) {
	var (
		vars map[string]string = mux.Vars(req)
		uid  string            = vars["uid"]
		ctx  context.Context   = req.Context()
	)

	batchStat, err := h.Usercases.IBatchStatUsecase.GetBatchStatById(ctx, uid)
	if err != nil {
		logger.Error("failed to get batchstat: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusInternalServerError, httpresponse.HTTPStatusInternalServerErrorMessage, nil)
		return
	}

	httpresponse.SendJSONResponse(resp, http.StatusOK, httpresponse.HTTPStatusOKMessage, batchStat)
}

// ApiMarkCompletedBatchStatHandler godoc
// @Summary      Mark completed BatchStat
// @Description  Mark completed a BatchStat by uid
// @Tags         apisql
// @Accept       json
// @Produce      json
// @Param        uid path string  true  "uid"
// @Success      200  {object}  httpresponse.HTTPResp
// @Failure      500  {object}  httpresponse.HTTPResp
// @Security     BasicAuth
// @Router       /api-gateway-sql/batchstats/{uid}/completed [get]
func (h *HTTPHandler) ApiMarkCompletedBatchStatHandler(resp http.ResponseWriter, req *http.Request) {
	var (
		vars map[string]string = mux.Vars(req)
		uid  string            = vars["uid"]
		ctx  context.Context   = req.Context()
		err  error
	)

	err = h.Usercases.IBatchStatUsecase.MarkCompletedBatchStat(ctx, uid)
	if err != nil {
		logger.Error("failed to mark compled a batchstat: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusInternalServerError, httpresponse.HTTPStatusInternalServerErrorMessage, nil)
		return
	}

	httpresponse.SendJSONResponse(resp, http.StatusOK, httpresponse.HTTPStatusOKMessage, nil)
}

// ApiListBlocksByBatchStatsHandler godoc
// @Summary      Get Blocks result
// @Description  Get a paginated list of Blocks By a BatchStat
// @Tags         apisql
// @Accept       json
// @Produce      json
// @Param        uid path string  true  "uid"
// @Param        page_num query  int  false  "Page number" default(1)
// @Param        page_size query int  false  "Page size" default(20)
// @Success      200  {object}  httpresponse.HTTPResp
// @Failure      400  {object}  httpresponse.HTTPResp
// @Failure      500  {object}  httpresponse.HTTPResp
// @Security     BasicAuth
// @Router       /api-gateway-sql/batchstats/{uid}/blocks [get]
func (h *HTTPHandler) ApiListBlocksByBatchStatHandler(resp http.ResponseWriter, req *http.Request) {
	var (
		vars     map[string]string = mux.Vars(req)
		uid      string            = vars["uid"]
		queries  url.Values        = req.URL.Query()
		pageNum  int
		pageSize int
		ctx      context.Context = req.Context()
		err      error
	)

	pageNum, err = strconv.Atoi(queries.Get("page_num"))
	if err != nil {
		logger.Error("failed to handle page_num param query: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusBadRequest, "Unable to handle page_num", nil)
		return
	}
	pageSize, err = strconv.Atoi(queries.Get("page_size"))
	if err != nil {
		logger.Error("failed to handle page_size param query: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusBadRequest, "Unable to handle page_size", nil)
		return
	}

	pageRequest := paginator.NewPageRequest(pageNum, pageSize)
	var blocksPaginatorResponse *paginator.PageResponse

	blocksPaginatorResponse, err = h.Usercases.IBlockUsecase.ListBlocksByBatchStat(ctx, uid, pageRequest)
	if err != nil {
		logger.Error("failed to list blocks: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusInternalServerError, httpresponse.HTTPStatusInternalServerErrorMessage, nil)
		return
	}

	httpresponse.SendJSONResponse(resp, http.StatusOK, httpresponse.HTTPStatusOKMessage, blocksPaginatorResponse)
}

// ApiGetBlockHandler godoc
// @Summary      Get Block result
// @Description  Get a Block by uid
// @Tags         apisql
// @Accept       json
// @Produce      json
// @Param        uid path string  true  "uid"
// @Success      200  {object}  httpresponse.HTTPResp
// @Failure      500  {object}  httpresponse.HTTPResp
// @Security     BasicAuth
// @Router       /api-gateway-sql/blocks/{uid} [get]
func (h *HTTPHandler) ApiGetBlockHandler(resp http.ResponseWriter, req *http.Request) {
	var (
		vars map[string]string = mux.Vars(req)
		uid  string            = vars["uid"]
		ctx  context.Context   = req.Context()
	)

	block, err := h.Usercases.IBlockUsecase.GetBlockById(ctx, uid)
	if err != nil {
		logger.Error("failed to get block: %s", err.Error())
		httpresponse.SendJSONResponse(resp, http.StatusInternalServerError, httpresponse.HTTPStatusInternalServerErrorMessage, nil)
		return
	}

	httpresponse.SendJSONResponse(resp, http.StatusOK, httpresponse.HTTPStatusOKMessage, block)
}
