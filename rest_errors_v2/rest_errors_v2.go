package rest_errors_v2

import (
	"fmt"
	"net/http"
)

type RestErr interface {
	Status() int     // HTTP status code
	Title() string   // A string representation of the Status Code
	Path() string    // Only used for Logging: The path of the error. Ex: "controller/controllerfunc/service/servicefunc/dbclient/dblientfunc"
	WrapPath(string) // Only used for Logging: Wrapper func to keep track of the path of the error
	Code() string    // Only used for Logging: Raw error code
	Message() string // Only used for Logging: Raw error message not returned to the client

	Error() string // string representation of a restErr
}

type restErr struct {
	ErrStatus  int    `json:"status"`          // HTTP Status Code
	ErrTitle   string `json:"title"`           // A string representation of the Status Code
	ErrCause   string `json:"cause,omitempty"` // The cause of the error, can be empty
	ErrPath    string `json:"-"`               // Only used for Logging: The path of the error. Ex: "controller/controllerfunc/service/servicefunc/dbclient/dblientfunc"
	ErrMessage string `json:"-"`               // Only used for Logging: Raw error message returned by a DB, another Servive or whatever
	ErrCode    string `json:"-"`               // Only used for Logging: Raw error code from the DB or another service
}

func (e *restErr) Error() string {
	return fmt.Sprintf("status: %d, title: %s, cause: %s - path: %s - message: %s",
		e.ErrStatus, e.ErrTitle, e.ErrCause, e.ErrPath, e.ErrMessage)
}

func (e *restErr) Status() int {
	return e.ErrStatus
}

func (e *restErr) Title() string {
	return e.ErrTitle
}

func (e *restErr) Message() string {
	return e.ErrMessage
}

func (e *restErr) Path() string {
	return e.ErrPath
}

func (e *restErr) Code() string {
	return e.ErrCode
}

func (e *restErr) WrapPath(path string) {
	e.ErrPath = fmt.Sprintf("%s%s", path, e.ErrPath)
}

// constructors
func NewInternalServerError(path string, code string, msg string) RestErr {
	result := &restErr{
		ErrStatus:  http.StatusInternalServerError,
		ErrTitle:   "internal_server_error",
		ErrMessage: msg,
		ErrCode:    code,
		ErrPath:    path,
	}
	return result
}

func NewBadRequestError(path string, message string, cause string) RestErr {
	return &restErr{
		ErrStatus:  http.StatusBadRequest,
		ErrTitle:   "bad_request",
		ErrMessage: message,
		ErrCause:   cause,
		ErrPath:    path,
	}
}
