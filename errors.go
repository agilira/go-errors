// errors.go: Core error types and constructors for the go-errors AGILira library
//
// Copyright (c) 2025 AGILira - A. Giordano
// Series: an AGLIra library
// SPDX-License-Identifier: MPL-2.0

package errors

import (
	"time"

	"github.com/agilira/go-timecache"
)

// ErrorCode represents a custom error code that can be used to categorize and identify specific types of errors.
// Error codes should be defined as constants in your application for consistency.
type ErrorCode string

// Predefined severity levels for consistent error classification.
// These constants can be used with WithSeverity() method for type safety.
const (
	SeverityCritical = "critical" // System failures, data corruption, security breaches
	SeverityError    = "error"    // Standard errors that prevent operation completion
	SeverityWarning  = "warning"  // Issues that don't prevent operation but need attention
	SeverityInfo     = "info"     // Informational messages for debugging/audit trails
)

// DefaultErrorCode is used when an empty or invalid ErrorCode is provided to constructors.
const DefaultErrorCode ErrorCode = "UNKNOWN_ERROR"

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
// The error will have a timestamp set to the current time and default severity of SeverityError.
// If code is empty or whitespace-only, DefaultErrorCode will be used instead.
//
// Example:
//
//	const ErrCodeValidation ErrorCode = "VALIDATION_ERROR"
//	err := New(ErrCodeValidation, "Username is required")
//	fmt.Println(err.Error()) // Output: [VALIDATION_ERROR]: Username is required
func New(code ErrorCode, message string) *Error {
	if !validateErrorCode(code) {
		code = DefaultErrorCode
	}
	return &Error{
		Code:      code,
		Message:   message,
		Timestamp: timecache.CachedTime(),
		Severity:  SeverityError,
		Context:   make(map[string]interface{}),
	}
}

// NewWithField creates a new structured error with the given code, message, field, and value.
// This is useful for validation errors where you need to specify which field caused the error.
// If code is empty or whitespace-only, DefaultErrorCode will be used instead.
//
// Example:
//
//	err := NewWithField("VALIDATION_ERROR", "Invalid email format", "email", "invalid@")
//	fmt.Printf("Field: %s, Value: %s\n", err.Field, err.Value)
//	// Output: Field: email, Value: invalid@
func NewWithField(code ErrorCode, message, field, value string) *Error {
	if !validateErrorCode(code) {
		code = DefaultErrorCode
	}
	return &Error{
		Code:      code,
		Message:   message,
		Field:     field,
		Value:     value,
		Timestamp: timecache.CachedTime(),
		Severity:  SeverityError,
		Context:   make(map[string]interface{}),
	}
}

// NewWithContext creates a new structured error with the given code, message, and context map.
// The context map allows you to attach additional metadata to the error for debugging purposes.
// If code is empty or whitespace-only, DefaultErrorCode will be used instead.
func NewWithContext(code ErrorCode, message string, context map[string]interface{}) *Error {
	if !validateErrorCode(code) {
		code = DefaultErrorCode
	}
	return &Error{
		Code:      code,
		Message:   message,
		Timestamp: timecache.CachedTime(),
		Severity:  SeverityError,
		Context:   context,
	}
}
