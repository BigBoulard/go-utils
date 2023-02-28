package rest_errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-resty/resty/v2"
)

type RestErr interface {
	Message() string
	Status() int
	Error() string
	Cause() string
}

type restErr struct {
	ErrMessage string `json:"message"`
	ErrStatus  int    `json:"status"`
	ErrError   string `json:"error"`
	ErrCause   string `json:"cause"`
}

func (e restErr) Error() string {
	return fmt.Sprintf("message: %s - status: %d - error: %s - causes: %v",
		e.ErrMessage, e.ErrStatus, e.ErrError, e.ErrCause)
}

func (e restErr) Message() string {
	return e.ErrMessage
}

func (e restErr) Status() int {
	return e.ErrStatus
}

func (e restErr) Cause() string {
	return e.ErrCause
}

func NewRestError(message string, status int, err string, cause string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  status,
		ErrError:   err,
		ErrCause:   cause,
	}
}

func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr restErr
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return apiErr, nil
}

func NewBadRequestError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusBadRequest,
		ErrError:   "bad_request",
	}
}

func NewServiceUnavailableError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusServiceUnavailable,
		ErrError:   "service_unavailable",
	}
}

func NewNotFoundError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusNotFound,
		ErrError:   "not_found",
	}
}

func NewGoneError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusGone,
		ErrError:   "gone",
	}
}

func NewUnauthorizedError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusUnauthorized,
		ErrError:   "unauthorized",
	}
}

func NewConflictError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusConflict,
		ErrError:   "conflict",
	}
}

func NewUnprocessableEntityError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusUnprocessableEntity,
		ErrError:   "Unprocessable Entity",
	}
}

func NewInternalServerError(message string, err error) RestErr {
	result := restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   "internal_server_error",
	}
	if err != nil {
		result.ErrCause = err.Error()
	}
	return result
}

func CheckRestError(err error, resp *resty.Response, origin string) RestErr {
	if err != nil {
		return NewServiceUnavailableError(
			fmt.Sprintf("%s:%s", origin, err.Error()),
		)
	}

	if resp.IsError() { // Network Error
		println("ISErrror")
		spew.Dump(resp)
		return NewInternalServerError(
			fmt.Sprintf("%s:%s", origin, resp.Error()),
			nil,
		)
	}

	if resp.StatusCode() > 399 {
		restErr := &restErr{}
		unmarshalErr := json.Unmarshal(resp.Body(), &restErr)
		if unmarshalErr != nil {
			spew.Dump(resp)
			println(resp.StatusCode())
			return NewInternalServerError(
				fmt.Sprintf("%s - Unmarshal error: %s", origin, unmarshalErr.Error()),
				unmarshalErr,
			)
		}
		spew.Dump(resp)
		return restErr
	}

	return nil
}
