// errors.go: Package errors provides reusable, structured error handling for Go applications.
//
// Copyright (c) 2025 AGILira
// Series: an AGLIra fragment
// SPDX-License-Identifier: MPL-2.0

package errors

import (
	"time"
)

// ErrorCode represents a custom error code (user-defined)
type ErrorCode string

// Error represents a structured error with context and cause
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

// New creates a new error with code and message
func New(code ErrorCode, message string) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Severity:  "error",
		Context:   make(map[string]interface{}),
	}
}

// NewWithField creates a new error with field and value
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

// NewWithContext creates a new error with additional context
func NewWithContext(code ErrorCode, message string, context map[string]interface{}) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Severity:  "error",
		Context:   context,
	}
}
