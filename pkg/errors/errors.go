package custom_errors

import (
	"errors"
	"fmt"
	"net/http"
)

type HTTPError interface {
	error
	StatusCode() int
	Unwrap() error
}

type BaseError struct {
	Version string `json:"version"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"cause,omitempty"`
}

func unwrapCause(cause []error) error {
	if len(cause) > 0 {
		return cause[0]
	}
	return nil
}

func (e *BaseError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *BaseError) StatusCode() int {
	return e.Code
}

func (e *BaseError) Unwrap() error {
	return e.Cause
}

type UnauthorizedError struct {
	BaseError
}

func NewUnauthorizedError(version string, message string, cause ...error) *UnauthorizedError {
	return &UnauthorizedError{
		BaseError: BaseError{
			Version: version,
			Code:    http.StatusUnauthorized,
			Message: message,
			Cause:   unwrapCause(cause),
		},
	}
}

type NotFoundError struct {
	BaseError
}

func NewNotFoundError(version string, message string, cause ...error) *NotFoundError {
	return &NotFoundError{
		BaseError: BaseError{
			Version: version,
			Code:    http.StatusNotFound,
			Message: message,
			Cause:   unwrapCause(cause),
		},
	}
}

func NewBadRequestError(version string, message string, cause ...error) *BadRequestError {

	return &BadRequestError{
		BaseError: BaseError{
			Version: version,
			Code:    http.StatusBadRequest,
			Message: message,
			Cause:   unwrapCause(cause),
		},
	}
}

type BadRequestError struct {
	BaseError
}

func NewInternalError(version string, message string, cause ...error) *InternalError {

	return &InternalError{
		BaseError: BaseError{
			Version: version,
			Code:    http.StatusInternalServerError,
			Message: message,
			Cause:   unwrapCause(cause),
		},
	}
}

type InternalError struct {
	BaseError
}

func NewForbiddenError(version string, message string, cause ...error) *ForbiddenError {

	return &ForbiddenError{
		BaseError: BaseError{
			Version: version,
			Code:    http.StatusForbidden,
			Message: message,
			Cause:   unwrapCause(cause),
		},
	}
}

type ForbiddenError struct {
	BaseError
}

func NewConflictError(version string, message string, cause ...error) *ConflictError {
	return &ConflictError{
		BaseError: BaseError{
			Version: version,
			Code:    http.StatusConflict,
			Message: message,
			Cause:   unwrapCause(cause),
		},
	}
}

type ConflictError struct {
	BaseError
}

func NewServiceUnavailableError(version string, message string, cause ...error) *ServiceUnavailableError {
	return &ServiceUnavailableError{
		BaseError: BaseError{
			Version: version,
			Code:    http.StatusServiceUnavailable,
			Message: message,
			Cause:   unwrapCause(cause),
		},
	}
}

type ServiceUnavailableError struct {
	BaseError
}

func IsConflict(err error) bool {
	if err == nil {
		return false
	}

	target := &ConflictError{}
	return errors.As(err, &target)
}
