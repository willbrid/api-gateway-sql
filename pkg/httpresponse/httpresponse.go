package httpresponse

import (
	"encoding/json"
	"net/http"
)

const (
	HTTPStatusOKMessage                  = "OK"
	HTTPStatusInternalServerErrorMessage = "Internal Server Error"
)

type HTTPResp struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"status ok"`
	Data    any    `json:"data,omitempty"`
}

func SendJSONResponse(resp http.ResponseWriter, status int, message string, data any) error {
	response := HTTPResp{
		Code:    status,
		Message: message,
		Data:    data,
	}

	resp.Header().Set("Content-Type", "application/json")

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return err
	}

	resp.WriteHeader(status)
	resp.Write(jsonResponse)

	return nil
}
