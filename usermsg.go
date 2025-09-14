// usermsg.go: Returns the user-friendly messages to the users
//
// Copyright (c) 2025 AGILira - A. Giordano
// Series: an AGLIra library
// SPDX-License-Identifier: MPL-2.0

package errors

// WithUserMessage sets a user-friendly message on the error and returns the error for chaining.
// This message should be safe to display to end users without exposing technical details.
//
// Example:
//
//	err := New("DB_CONNECTION_ERROR", "Database connection timeout").
//		WithUserMessage("We're experiencing technical difficulties. Please try again later.")
func (e *Error) WithUserMessage(msg string) *Error {
	e.UserMsg = msg
	return e
}

// WithContext adds or updates context information on the error and returns the error for chaining.
// Context information is useful for debugging and logging purposes.
func (e *Error) WithContext(key string, value interface{}) *Error {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// AsRetryable marks the error as retryable and returns the error for chaining.
// This indicates that the operation that caused this error can be safely retried.
func (e *Error) AsRetryable() *Error {
	e.Retryable = true
	return e
}

// WithSeverity sets the severity level of the error and returns the error for chaining.
// Common severity levels include "error", "warning", "info", and "critical".
func (e *Error) WithSeverity(severity string) *Error {
	e.Severity = severity
	return e
}

// UserMessage returns the user-friendly message if set, otherwise falls back to the technical message.
// This implements the UserMessager interface.
func (e *Error) UserMessage() string {
	if e.UserMsg != "" {
		return e.UserMsg
	}
	return e.Message
}

// ErrorCode returns the error code.
// This implements the ErrorCoder interface.
func (e *Error) ErrorCode() ErrorCode {
	return e.Code
}

// IsRetryable returns whether the error is retryable.
// This implements the Retryable interface.
func (e *Error) IsRetryable() bool {
	return e.Retryable
}

// WithCriticalSeverity sets the error severity to critical and returns the error for chaining.
// Use this for system failures, data corruption, or security breaches.
func (e *Error) WithCriticalSeverity() *Error {
	return e.WithSeverity(SeverityCritical)
}

// WithWarningSeverity sets the error severity to warning and returns the error for chaining.
// Use this for issues that don't prevent operation but need attention.
func (e *Error) WithWarningSeverity() *Error {
	return e.WithSeverity(SeverityWarning)
}

// WithInfoSeverity sets the error severity to info and returns the error for chaining.
// Use this for informational messages for debugging or audit trails.
func (e *Error) WithInfoSeverity() *Error {
	return e.WithSeverity(SeverityInfo)
}
