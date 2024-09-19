package errors

import (
	goerrors "errors"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

var (
	ErrInternalServerError = goerrors.New("internal server error")

	ErrNotFound = goerrors.New("your requested item is not found")
)

type FieldError struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
	Tag   string `json:"tag"`
}

type ValidationError struct {
	Errors []FieldError `json:"errors"`
}

func (ve ValidationError) Error() string {
	res := "ValidationError ("
	for _, e := range ve.Errors {
		res += e.Msg + "; "
	}
	return strings.TrimSuffix(res, "; ") + ")"
}

// RecoveredError is an error that indicates it should not stop the program.
// Just log the error if it's occured, don't panic.
type RecoveredError struct {
	ActualErr error
}

func (e RecoveredError) Error() string {
	if e.ActualErr != nil {
		return e.ActualErr.Error()
	}
	return ""
}

type ResponseError struct {
	Code    int
	ErrCode string
	Message string
	Err     error // original error
}

func (r ResponseError) Error() string {
	return fmt.Sprintf("%d: %s", r.Code, r.Message)
}

func BadRequest(msg string) ResponseError {
	return ResponseError{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
}

type EntityNotFoundError struct {
	Message string
}

func (r EntityNotFoundError) Error() string {
	return r.Message
}

func (r EntityNotFoundError) Is(err error) bool {
	var vErr EntityNotFoundError
	return errors.As(err, &vErr)
}

func NewEntityNotFoundError(entity string) EntityNotFoundError {
	return EntityNotFoundError{
		Message: fmt.Sprintf("%s not found", entity),
	}
}

func IsEntityNotFoundErr(err error) bool {
	var enErr EntityNotFoundError
	if goerrors.As(err, &enErr) {
		return true
	}
	return false
}

func EntityNotFoundErrTo404ResponseErr(err error, message ...string) error {
	if err == nil {
		return nil
	}
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}

	var enErr EntityNotFoundError
	if goerrors.As(err, &enErr) {
		if msg == "" {
			msg = enErr.Message
		}
		return ResponseError{
			Code:    http.StatusNotFound,
			Message: msg,
			Err:     err,
		}
	}

	return err
}

func MappingError(err error, mp map[error]error) error {
	for kerr, verr := range mp {
		if errors.Is(err, kerr) {
			return verr
		}
	}
	return err
}
