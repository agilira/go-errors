// Package errors provides structured, contextual error handling for Go applications.
//
// # Overview
//
// This package offers a comprehensive error handling system designed for modern Go applications.
// It provides structured error types with rich metadata, stack traces, user-friendly messages,
// and JSON serialization capabilities. Built for microservices and API development, it enables
// consistent error handling across your entire application stack.
//
// # Key Features
//
// • Structured Error Types: Errors include codes, messages, context, timestamps, and severity levels
// • Stack Trace Support: Optional lightweight stack traces for debugging
// • User-Friendly Messages: Separate technical and user-facing error messages
// • JSON Serialization: Built-in JSON marshaling for API responses and logging
// • Retry Logic: Built-in support for retryable errors
// • Interface-Based: Type-safe error handling through well-defined interfaces
// • Zero Dependencies: Uses only Go standard library
// • High Performance: Minimal overhead with efficient memory usage
//
// # Quick Start
//
// Create a new error with a custom code:
//
//	const ErrCodeValidation = "VALIDATION_ERROR"
//	err := errors.New(ErrCodeValidation, "Username is required")
//
// Add user-friendly message and context:
//
//	err = err.WithUserMessage("Please enter a username").
//		WithContext("field", "username").
//		WithSeverity("warning")
//
// Wrap existing errors with additional context:
//
//	if err := someOperation(); err != nil {
//		return errors.Wrap(err, ErrCodeOperation, "Failed to process request")
//	}
//
// Check error codes programmatically:
//
//	if errors.HasCode(err, ErrCodeValidation) {
//		// Handle validation error
//	}
//
// # JSON Serialization
//
// Errors automatically serialize to JSON with all metadata:
//
//	jsonData, _ := json.Marshal(err)
//	// Output: {"code":"VALIDATION_ERROR","message":"Username is required",...}
//
// # Interface-Based Error Handling
//
// Use interfaces for type-safe error handling:
//
//	var coder errors.ErrorCoder = err
//	code := coder.ErrorCode()
//
//	var retry errors.Retryable = err
//	if retry.IsRetryable() {
//		// Implement retry logic
//	}
//
//	var um errors.UserMessager = err
//	userMsg := um.UserMessage()
//
// # Best Practices
//
// 1. Define error codes as constants in your application
// 2. Use WithUserMessage() for errors that will be displayed to users
// 3. Use WithContext() to add debugging information
// 4. Use Wrap() to add context to errors from lower-level functions
// 5. Use AsRetryable() for transient errors that can be retried
// 6. Use different severity levels for different types of errors
//
// # Error Severity Levels
//
// • "error": Standard application errors (default)
// • "warning": Non-critical issues that should be noted
// • "info": Informational messages
// • "critical": Severe errors that require immediate attention
//
// # Integration Examples
//
// REST API Handler:
//
//	func handleRequest(w http.ResponseWriter, r *http.Request) {
//		err := processRequest(r)
//		if err != nil {
//			apiErr, ok := err.(*errors.Error)
//			if ok {
//				http.Error(w, apiErr.UserMessage(), getStatusCode(apiErr))
//				log.Error("API Error", "error", apiErr.Error(), "code", apiErr.ErrorCode())
//			} else {
//				http.Error(w, "Internal server error", 500)
//			}
//		}
//	}
//
// Database Operations:
//
//	func saveUser(user *User) error {
//		if err := db.Save(user); err != nil {
//			return errors.Wrap(err, ErrCodeDatabase, "Failed to save user").
//				WithContext("user_id", user.ID).
//				WithUserMessage("Unable to save your information. Please try again.")
//		}
//		return nil
//	}
//
// Validation:
//
//	func validateUser(user *User) error {
//		if user.Email == "" {
//			return errors.NewWithField(ErrCodeValidation, "Email is required", "email", user.Email).
//				WithUserMessage("Please provide a valid email address")
//		}
//		return nil
//	}
//
// # Testing
//
// Use table-driven tests for error scenarios:
//
//	func TestValidation(t *testing.T) {
//		tests := []struct {
//			name    string
//			input   string
//			wantErr bool
//			wantCode errors.ErrorCode
//		}{
//			{"valid email", "user@example.com", false, ""},
//			{"empty email", "", true, ErrCodeValidation},
//		}
//
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				err := validateEmail(tt.input)
//				if tt.wantErr {
//					assert.True(t, errors.HasCode(err, tt.wantCode))
//				} else {
//					assert.NoError(t, err)
//				}
//			})
//		}
//	}
//
// # Performance Considerations
//
// • Stack traces are only captured when using Wrap() or explicitly requested
// • JSON marshaling is optimized for common use cases
// • Memory usage is minimal with efficient struct layout
// • No reflection is used in hot paths
//
// # Migration from Standard Errors
//
// Replace standard error creation:
//
//	// Before
//	return errors.New("validation failed")
//
//	// After
//	return errors.New(ErrCodeValidation, "validation failed")
//
// Replace error wrapping:
//
//	// Before
//	return fmt.Errorf("operation failed: %w", err)
//
//	// After
//	return errors.Wrap(err, ErrCodeOperation, "operation failed")
//
// Copyright (c) 2025 AGILira
// Series: an AGLIra library
// SPDX-License-Identifier: MPL-2.0
package errors
