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
