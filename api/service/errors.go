package service

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ErrorResponse) Error400(err error) {
	e.Code = http.StatusBadRequest
	e.Message = "Bad Request: " + err.Error()
	log.Error(err)
}

func (e *ErrorResponse) Error500(err error) {
	e.Code = http.StatusInternalServerError
	e.Message = "Internal Server Error: " + err.Error()
	log.Error(err)
}

func (e *ErrorResponse) UuidParseError(err error) {
	e.Code = http.StatusBadRequest
	e.Message = "ID read error: " + err.Error()
	log.Error(err)
}

func (e *ErrorResponse) ReadBodyError(err error) {
	e.Code = http.StatusBadRequest
	e.Message = "Read body error: " + err.Error()
	log.Error(err)
}

func (e *ErrorResponse) DeserializeError(err error) {
	e.Code = http.StatusInternalServerError
	e.Message = "Deserialize error: " + err.Error()
	log.Error(err)
}

func (e *ErrorResponse) SerializeError(err error) {
	e.Code = http.StatusInternalServerError
	e.Message = "Serialize error: " + err.Error()
	log.Error(err)
}

func (e *ErrorResponse) ValidationError(err error) {
	e.Code = http.StatusBadRequest
	e.Message = "Validation error: " + err.Error()
	log.Error(err)
}

func (e *ErrorResponse) DBError(err error) {
	e.Code = http.StatusInternalServerError
	e.Message = "Database error: " + err.Error()
	log.Error(err)
}

func (e *ErrorResponse) DBNotFound(err error) {
	e.Code = http.StatusNotFound
	e.Message = "Record not found: " + err.Error()
	log.Error(err)
}

func (e *ErrorResponse) TaskNotStartedError() {
	e.Code = http.StatusBadRequest
	e.Message = "Task not started."
	log.Error("Task not started.")
}

func (e *ErrorResponse) TaskIsAlreadyStartedError() {
	e.Code = http.StatusBadRequest
	e.Message = "Task is already started."
	log.Error("Task is already started.")
}

func (e *ErrorResponse) TaskIsAlreadyFinishedError() {
	e.Code = http.StatusBadRequest
	e.Message = "Task is already finished."
	log.Error("Task is already finished.")
}

func (e *ErrorResponse) ExternalAPIError(err error) {
	e.Code = http.StatusInternalServerError
	e.Message = "External API Error: " + err.Error()
	log.Error(err)
}
