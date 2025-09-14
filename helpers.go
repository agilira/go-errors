// helpers.go: Helpers functions for the go-errors AGILira library
//
// Copyright (c) 2025 AGILira - A. Giordano
// Series: an AGLIra library
// SPDX-License-Identifier: MPL-2.0

package errors

import (
	"errors"
	"fmt"

	"github.com/agilira/go-timecache"
)

// Wrap wraps an existing error with a new code and message, capturing the current stack trace.
// This is useful for adding context to errors that occur deeper in the call stack.
// If code is empty or whitespace-only, DefaultErrorCode will be used instead.
//
// Example:
//
//	if err := someOperation(); err != nil {
//		return Wrap(err, "OPERATION_FAILED", "Failed to process user data")
//	}
func Wrap(err error, code ErrorCode, message string) *Error {
	if !validateErrorCode(code) {
		code = DefaultErrorCode
	}
	return &Error{
		Code:      code,
		Message:   message,
		Timestamp: timecache.CachedTime(),
		Severity:  SeverityError,
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
//
// Example:
//
//	if HasCode(err, "VALIDATION_ERROR") {
//		// Handle validation-specific error
//		log.Warning("Validation failed", "error", err)
//	}
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

// validateErrorCode checks if an ErrorCode is valid (non-empty).
// Returns true if the code is valid, false if empty or whitespace-only.
func validateErrorCode(code ErrorCode) bool {
	s := string(code)
	if len(s) == 0 {
		return false
	}
	// Check if it's only whitespace
	for _, r := range s {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			return true
		}
	}
	return false
}
