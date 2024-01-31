package apperror

import (
	"fmt"
	"net/http"

	"github.com/sangianpatrick/tm-user/pkg/appstatus"
)

type Error struct {
	Err            error
	Status         appstatus.Status
	HTTPStatusCode int
}

func SetError(err error, status appstatus.Status, httpStatusCode int) *Error {
	if httpStatusCode < 400 {
		httpStatusCode = 500
		status = "InternalServerError"
	}
	if err == nil {
		err = fmt.Errorf("%d: %s", httpStatusCode, status)
	}

	return &Error{
		Err:            err,
		Status:         status,
		HTTPStatusCode: httpStatusCode,
	}
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func Destruct(err error) *Error {
	e, ok := err.(*Error)
	if ok {
		return e
	}

	e = new(Error)
	e.Err = err
	e.Status = appstatus.InternalServerError
	e.HTTPStatusCode = http.StatusInternalServerError

	return e
}
