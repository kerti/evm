package response

import (
	"encoding/json"
	"net/http"

	"github.com/kerti/evm/01-machine-gun/util/failure"
)

var failureStatusMap = map[failure.Code]int{
	failure.CodeBadRequest:            http.StatusBadRequest,
	failure.CodeUnauthorized:          http.StatusUnauthorized,
	failure.CodeInternalError:         http.StatusInternalServerError,
	failure.CodeUnimplemented:         http.StatusNotImplemented,
	failure.CodeEntityNotFound:        http.StatusNotFound,
	failure.CodeOperationNotPermitted: http.StatusConflict,
}

// BaseResponse is the base object of all responses
type BaseResponse struct {
	Data    *interface{} `json:"data,omitempty"`
	Error   *string      `json:"error,omitempty"`
	Message *string      `json:"message,omitempty"`
}

// RespondWithNoContent sends a response without any content
func RespondWithNoContent(w http.ResponseWriter) {
	respond(w, http.StatusNoContent, nil)
}

// RespondWithMessage sends a response with a simple text message
func RespondWithMessage(w http.ResponseWriter, code int, message string) {
	respond(w, code, BaseResponse{Message: &message})
}

// RespondWithJSON sends a response containing a JSON object
func RespondWithJSON(w http.ResponseWriter, code int, jsonPayload interface{}) {
	respond(w, code, BaseResponse{Data: &jsonPayload})
}

// RespondWithError sends a response with an error message
func RespondWithError(w http.ResponseWriter, err error) {
	code := failure.GetCode(err)
	status, ok := failureStatusMap[code]
	if !ok {
		status = http.StatusInternalServerError
	}
	errMsg := err.Error()
	respond(w, status, BaseResponse{Error: &errMsg})
}

// RespondWithPreparingShutdown sends a default response for when the server is preparing to shut down
func RespondWithPreparingShutdown(w http.ResponseWriter) {
	RespondWithMessage(w, http.StatusServiceUnavailable, "SERVER PREPARING TO SHUT DOWN")
}

// RespondWithUnhealthy sends a default response for when the server is unhealthy
func RespondWithUnhealthy(w http.ResponseWriter) {
	RespondWithMessage(w, http.StatusServiceUnavailable, "SERVER UNHEALTHY")
}

func respond(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
