package handler

import (
	"fmt"
	"strings"
)

// error message code
var (
	ErrFormatInvalid = "FORMAT_INVALID"
	ErrFailed        = "FAILED"
	ErrNotFound      = "NOT_FOUND"
	ErrValueInvalid  = "VALUE_INVALID"
	ErrRequired      = "REQUIRED"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e FieldError) String() string {
	return fmt.Sprintf("%v: %v", e.Field, e.Message)
}

type FieldErrors []FieldError

func (es FieldErrors) Error() string {
	var buf = make([]string, len(es))
	for i, e := range es {
		buf[i] = e.String()
	}
	return strings.Join(buf, ", ")
}
