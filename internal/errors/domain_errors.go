package errors

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type DomainError struct {
	msg     string
	Message string
	Code    string
}

func (r DomainError) Error() string {
	return fmt.Sprintf("%s (%s)", r.Code, r.Message)
}

func (r DomainError) Is(err error) bool {
	var vErr DomainError
	return errors.As(err, &vErr) && r.Code == vErr.Code && r.msg == vErr.msg
}

func NewDomainError(code string, message string) DomainError {
	if code != "" {
		code = "xyz-" + code
	}

	return DomainError{
		msg:     message,
		Message: message,
		Code:    code,
	}
}

var DmErrPermissionDenied = NewDomainError("", "you don't have permission to access this resource")

var ErrFieldRequired = func(field string) DomainError {
	return NewDomainError("general-0001", fmt.Sprintf("%s is required", field))
}

func FormatDomainMessage(err DomainError, replaces map[string]string) DomainError {
	for k, v := range replaces {
		err.Message = strings.ReplaceAll(err.Message, ":"+k, v)
	}
	return err
}

func DomainErrorToResponseError(err error) error {
	if err == nil {
		return nil
	}

	var dmErr DomainError
	if errors.As(err, &dmErr) {
		httpCode := http.StatusBadRequest
		if errors.Is(dmErr, DmErrPermissionDenied) {
			httpCode = http.StatusForbidden
		}
		err = ResponseError{
			Code:    httpCode,
			ErrCode: dmErr.Code,
			Message: dmErr.Message,
			Err:     err,
		}
	}

	return err
}
