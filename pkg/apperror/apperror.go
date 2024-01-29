package apperror

import (
	"fmt"

	"github.com/sangianpatrick/tm-user/pkg/appstatus"
)

type Error struct {
	err            error
	status         appstatus.Status
	httpStatusCode int
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
		err:            err,
		status:         status,
		httpStatusCode: httpStatusCode,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s %s", e.httpStatusCode, e.status, e.err.Error())
}
