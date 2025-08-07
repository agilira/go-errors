// interfaces.go: Interfaces for the go-errors AGILira library
//
// Copyright (c) 2025 AGILira
// Series: an AGLIra library
// SPDX-License-Identifier: MPL-2.0

package errors

// ErrorCoder allows extracting a code from an error.
type ErrorCoder interface {
	ErrorCode() ErrorCode
}

// Retryable marks an error as retryable.
type Retryable interface {
	IsRetryable() bool
}

// UserMessager allows extracting a user-friendly message from an error.
type UserMessager interface {
	UserMessage() string
}
