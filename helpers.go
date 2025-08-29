// helpers.go: Helpers functions for the go-errors AGILira library
//
// Copyright (c) 2025 AGILira
// Series: an AGLIra library
// SPDX-License-Identifier: MPL-2.0

package errors

import (
	"errors"
	"fmt"
	"time"
)

// Wrap wraps an existing error with a new code and message, capturing the current stack trace.
// This is useful for adding context to errors that occur deeper in the call stack.
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
// It returns a formatted string containing the error code and message.
func (e *Error) Error() string {
	return fmt.Sprintf("[%s]: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause error, implementing the error wrapping interface.
func (e *Error) Unwrap() error {
	return e.Cause
}

// RootCause returns the original error in the error chain by unwrapping all nested errors.
// This is useful for finding the root cause of an error that has been wrapped multiple times.
func RootCause(err error) error {
	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			return err
		}
		err = unwrapped
	}
}

// HasCode checks if any error in the error chain has the given error code.
// This is useful for checking if a specific type of error occurred anywhere in the chain.
func HasCode(err error, code ErrorCode) bool {
	for err != nil {
		if ec, ok := err.(*Error); ok && ec.Code == code {
			return true
		}
		err = errors.Unwrap(err)
	}
	return false
}

// Is implements errors.Is compatibility for error comparison.
// It returns true if the target error has the same error code.
func (e *Error) Is(target error) bool {
	if target == nil {
		return false
	}
	if te, ok := target.(*Error); ok {
		return e.Code == te.Code
	}
	return false
}

// As implements errors.As compatibility for error type assertion.
// It delegates the check to the underlying Cause error.
// Note: It will not match the *Error instance itself, only errors in its cause chain.
func (e *Error) As(target interface{}) bool {
	return errors.As(e.Cause, target)
}
