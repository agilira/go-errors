// interfaces.go: Interfaces for the go-errors AGILira library
//
// Copyright (c) 2025 AGILira - A. Giordano
// Series: an AGLIra library
// SPDX-License-Identifier: MPL-2.0

package errors

// ErrorCoder allows extracting an error code from an error.
// This interface enables type-safe error code checking without type assertions.
type ErrorCoder interface {
	ErrorCode() ErrorCode
}

// Retryable indicates whether an error represents a condition that can be safely retried.
// This interface is useful for implementing retry logic in applications.
type Retryable interface {
	IsRetryable() bool
}

// UserMessager allows extracting a user-friendly message from an error.
// This interface enables displaying safe, non-technical messages to end users.
type UserMessager interface {
	UserMessage() string
}
