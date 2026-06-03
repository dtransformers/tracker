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

type ErrorResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func RespondError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var httpErr custom_errors.HTTPError
	if errors.As(err, &httpErr) {
		status := httpErr.StatusCode()
		c.JSON(status, ErrorResponse{
			Status: "error",
			Code:   status,
			Error:  httpErr.GetMessage(),
		})
		return
	}

	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Status: "error",
		Code:   http.StatusInternalServerError,
		Error:  "Internal server error",
	})
}

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
