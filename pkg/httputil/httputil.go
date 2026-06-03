package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	gojson "github.com/goccy/go-json"

	custom_errors "tracking/pkg/errors"
)

// ErrorResponse is the shared JSON error envelope used by all handlers.
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
	Status  int    `json:"status"`
	Success bool   `json:"success"`
	TraceID string `json:"trace_id,omitempty"`
	Version string `json:"version,omitempty"`
}

// RespondError converts a typed pkg/errors error into the correct HTTP response.
// It never exposes internal cause details (stack traces, DB errors, etc.) to the
// client. Unknown/internal errors are responded to with a generic 500 message.
func RespondError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var httpErr custom_errors.HTTPError
	if errors.As(err, &httpErr) {
		status := httpErr.StatusCode()
		c.JSON(status, ErrorResponse{
			Error:   httpErr.Error(),
			Status:  status,
			Success: false,
		})
		return
	}

	// Unknown error — return a generic 500 without leaking internals.
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error:   "Internal server error",
		Status:  http.StatusInternalServerError,
		Success: false,
	})
}

// FormatBindError produces a client-safe description of a JSON bind or validation error.
func FormatBindError(err error) string {
	if err == nil {
		return "invalid request body"
	}

	var unmarshalTypeErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeErr) {
		return fmt.Sprintf("Invalid data type for field: %s", unmarshalTypeErr.Field)
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return "Malformed JSON payload"
	}

	// goccy/go-json type error
	var goJsonTypeErr *gojson.UnmarshalTypeError
	if errors.As(err, &goJsonTypeErr) {
		return fmt.Sprintf("Invalid data type for field: %s", goJsonTypeErr.Field)
	}

	var goJsonSyntaxErr *gojson.SyntaxError
	if errors.As(err, &goJsonSyntaxErr) {
		return "Malformed JSON payload"
	}

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		if len(validationErrs) > 0 {
			return fmt.Sprintf("Validation failed on field: %s", validationErrs[0].Field())
		}
	}

	return "The provided payload format is invalid or contains mismatched data types."
}
