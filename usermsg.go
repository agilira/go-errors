// usermsg.go: Returns the user-friendly messages to the users
//
// Copyright (c) 2025 AGILira
// Series: an AGLIra library
// SPDX-License-Identifier: MPL-2.0

package errors

// WithUserMessage sets a user-friendly message on the error.
func (e *Error) WithUserMessage(msg string) *Error {
	e.UserMsg = msg
	return e
}

// WithContext adds or updates context information on the error.
func (e *Error) WithContext(key string, value interface{}) *Error {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// AsRetryable marks the error as retryable.
func (e *Error) AsRetryable() *Error {
	e.Retryable = true
	return e
}

// WithSeverity sets the severity level of the error.
func (e *Error) WithSeverity(severity string) *Error {
	e.Severity = severity
	return e
}

// UserMessage returns the user-friendly message if set, otherwise the technical message.
func (e *Error) UserMessage() string {
	if e.UserMsg != "" {
		return e.UserMsg
	}
	return e.Message
}

// ErrorCode returns the error code.
func (e *Error) ErrorCode() ErrorCode {
	return e.Code
}

// IsRetryable returns whether the error is retryable.
func (e *Error) IsRetryable() bool {
	return e.Retryable
}
