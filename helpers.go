// helpers.go: Helpers functions for the go-errors AGILira library
//
// Copyright (c) 2025 AGILira
// Series: an AGLIra fragment
// SPDX-License-Identifier: MPL-2.0

package errors

import (
	"errors"
	"fmt"
	"time"
)

// Wrap wraps an existing error with code and message, capturing stacktrace.
func Wrap(err error, code ErrorCode, message string) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Severity:  "error",
		Cause:     err,
		Context:   make(map[string]interface{}),
		Stack:     CaptureStacktrace(1),
	}
}

// Error implements the error interface for *Error.
func (e *Error) Error() string {
	return fmt.Sprintf("[%s]: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause.
func (e *Error) Unwrap() error {
	return e.Cause
}

// RootCause returns the original error in the chain.
func RootCause(err error) error {
	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			return err
		}
		err = unwrapped
	}
}

// HasCode checks if any error in the chain has the given code.
func HasCode(err error, code ErrorCode) bool {
	for err != nil {
		if ec, ok := err.(*Error); ok && ec.Code == code {
			return true
		}
		err = errors.Unwrap(err)
	}
	return false
}

// Is implements errors.Is compatibility.
func (e *Error) Is(target error) bool {
	if target == nil {
		return false
	}
	if te, ok := target.(*Error); ok {
		return e.Code == te.Code
	}
	return false
}

// As implements errors.As compatibility. It delegates the check to the underlying Cause.
// Note: It will not match the *Error instance itself, only errors in its cause chain.
func (e *Error) As(target interface{}) bool {
	return errors.As(e.Cause, target)
}
