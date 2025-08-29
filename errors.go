// errors.go: Core error types and constructors for the go-errors AGILira library
//
// Copyright (c) 2025 AGILira
// Series: an AGLIra library
// SPDX-License-Identifier: MPL-2.0

package errors

import (
	"time"
)

// ErrorCode represents a custom error code that can be used to categorize and identify specific types of errors.
// Error codes should be defined as constants in your application for consistency.
type ErrorCode string

// Error represents a structured error with comprehensive context and metadata.
// It includes error codes, messages, stack traces, user-friendly messages, and retry information.
type Error struct {
	Code      ErrorCode              `json:"code"`
	Message   string                 `json:"message"`
	Field     string                 `json:"field,omitempty"`
	Value     string                 `json:"value,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Cause     error                  `json:"cause,omitempty"`
	Severity  string                 `json:"severity"`
	Stack     *Stacktrace            `json:"stack,omitempty"`
	UserMsg   string                 `json:"user_msg,omitempty"`
	Retryable bool                   `json:"retryable,omitempty"`
}

// New creates a new structured error with the given code and message.
// The error will have a timestamp set to the current time and default severity of "error".
func New(code ErrorCode, message string) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Severity:  "error",
		Context:   make(map[string]interface{}),
	}
}

// NewWithField creates a new structured error with the given code, message, field, and value.
// This is useful for validation errors where you need to specify which field caused the error.
func NewWithField(code ErrorCode, message, field, value string) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Field:     field,
		Value:     value,
		Timestamp: time.Now(),
		Severity:  "error",
		Context:   make(map[string]interface{}),
	}
}

// NewWithContext creates a new structured error with the given code, message, and context map.
// The context map allows you to attach additional metadata to the error for debugging purposes.
func NewWithContext(code ErrorCode, message string, context map[string]interface{}) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Severity:  "error",
		Context:   context,
	}
}
