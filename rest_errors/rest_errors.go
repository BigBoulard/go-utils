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
	Causes() []interface{}
}

type restErr struct {
	ErrMessage string        `json:"message"`
	ErrStatus  int           `json:"status"`
	ErrError   string        `json:"error"`
	ErrCauses  []interface{} `json:"causes"`
}

func (e restErr) Error() string {
	return fmt.Sprintf("message: %s - status: %d - error: %s - causes: %v",
		e.ErrMessage, e.ErrStatus, e.ErrError, e.ErrCauses)
}

func (e restErr) Message() string {
	return e.ErrMessage
}

func (e restErr) Status() int {
	return e.ErrStatus
}

func (e restErr) Causes() []interface{} {
	return e.ErrCauses
}

func NewRestError(message string, status int, err string, causes []interface{}) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  status,
		ErrError:   err,
		ErrCauses:  causes,
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

func NewInternalServerError(message string, err error) RestErr {
	result := restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   "internal_server_error",
	}
	if err != nil {
		result.ErrCauses = append(result.ErrCauses, err.Error())
	}
	return result
}

func CheckRestError(err error, resp *resty.Response, origin string) RestErr {
	if err != nil { // TODO detail error scenarios here
		spew.Dump(err)
		return NewBadRequestError(
			fmt.Sprintf("%s:%s", origin, err.Error()),
		)
	}

	if resp.IsError() {
		restErr := &restErr{} // we assume that we only receive errors of type RestErr
		unmarshalErr := json.Unmarshal(resp.Body(), &restErr)
		if unmarshalErr != nil {
			return NewInternalServerError(
				fmt.Sprintf("%s - Unmarshal error: %s", origin, unmarshalErr.Error()),
				unmarshalErr,
			)
		}
		return restErr
	}

	return nil
}
