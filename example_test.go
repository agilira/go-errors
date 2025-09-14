// example_test.go: Examples for the go-errors AGILira library
//
// Copyright (c) 2025 AGILira - A. Giordano
// Series: an AGLIra library
// SPDX-License-Identifier: MPL-2.0

package errors_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/agilira/go-errors"
)

// Example error codes - define these as constants in your application
const (
	ErrCodeValidation = "VALIDATION_ERROR"
	ErrCodeDatabase   = "DATABASE_ERROR"
	ErrCodeNetwork    = "NETWORK_ERROR"
	ErrCodeAuth       = "AUTHENTICATION_ERROR"
)

// ExampleNew demonstrates basic error creation
func ExampleNew() {
	err := errors.New(ErrCodeValidation, "Username is required")
	fmt.Println(err.Error())
	fmt.Printf("Code: %s\n", err.ErrorCode())
	// Output:
	// [VALIDATION_ERROR]: Username is required
	// Code: VALIDATION_ERROR
}

// ExampleNewWithField demonstrates field-specific validation errors
func ExampleNewWithField() {
	err := errors.NewWithField(ErrCodeValidation, "Invalid email format", "email", "invalid@")
	fmt.Printf("Error: %s\n", err.Error())
	fmt.Printf("Field: %s, Value: %s\n", err.Field, err.Value)
	// Output:
	// Error: [VALIDATION_ERROR]: Invalid email format
	// Field: email, Value: invalid@
}

// ExampleError_WithUserMessage demonstrates user-friendly error messages
func ExampleError_WithUserMessage() {
	err := errors.New(ErrCodeDatabase, "Connection timeout after 30 seconds").
		WithUserMessage("We're experiencing technical difficulties. Please try again later.")

	fmt.Printf("Technical: %s\n", err.Message)
	fmt.Printf("User-friendly: %s\n", err.UserMessage())
	// Output:
	// Technical: Connection timeout after 30 seconds
	// User-friendly: We're experiencing technical difficulties. Please try again later.
}

// ExampleError_WithContext demonstrates adding debugging context
func ExampleError_WithContext() {
	err := errors.New(ErrCodeDatabase, "Query failed").
		WithContext("query", "SELECT * FROM users").
		WithContext("duration_ms", 5000).
		WithContext("connection_pool", "primary")

	fmt.Printf("Error: %s\n", err.Error())
	fmt.Printf("Query: %s\n", err.Context["query"])
	fmt.Printf("Duration: %v ms\n", err.Context["duration_ms"])
	// Output:
	// Error: [DATABASE_ERROR]: Query failed
	// Query: SELECT * FROM users
	// Duration: 5000 ms
}

// ExampleWrap demonstrates error wrapping with stack traces
func ExampleWrap() {
	// Simulate a low-level error
	originalErr := fmt.Errorf("connection refused")

	// Wrap with structured error
	wrappedErr := errors.Wrap(originalErr, ErrCodeNetwork, "Failed to connect to service")

	fmt.Printf("Wrapped: %s\n", wrappedErr.Error())
	fmt.Printf("Root cause: %s\n", errors.RootCause(wrappedErr).Error())
	// Output:
	// Wrapped: [NETWORK_ERROR]: Failed to connect to service
	// Root cause: connection refused
}

// ExampleError_severity demonstrates severity levels
func ExampleError_severity() {
	// Using predefined severity constants
	criticalErr := errors.New(ErrCodeDatabase, "Data corruption detected").
		WithCriticalSeverity()

	warningErr := errors.New(ErrCodeValidation, "Optional field missing").
		WithWarningSeverity()

	infoErr := errors.New("INFO_CODE", "Operation completed successfully").
		WithInfoSeverity()

	fmt.Printf("Critical: %s (severity: %s)\n", criticalErr.Message, criticalErr.Severity)
	fmt.Printf("Warning: %s (severity: %s)\n", warningErr.Message, warningErr.Severity)
	fmt.Printf("Info: %s (severity: %s)\n", infoErr.Message, infoErr.Severity)
	// Output:
	// Critical: Data corruption detected (severity: critical)
	// Warning: Optional field missing (severity: warning)
	// Info: Operation completed successfully (severity: info)
}

// ExampleError_AsRetryable demonstrates retry logic
func ExampleError_AsRetryable() {
	retryableErr := errors.New(ErrCodeNetwork, "Temporary service unavailable").
		AsRetryable().
		WithUserMessage("Service temporarily unavailable. Please try again.")

	if retryableErr.IsRetryable() {
		fmt.Println("This error can be retried")
		fmt.Printf("User message: %s\n", retryableErr.UserMessage())
	}
	// Output:
	// This error can be retried
	// User message: Service temporarily unavailable. Please try again.
}

// ExampleError_MarshalJSON demonstrates JSON serialization
func ExampleError_MarshalJSON() {
	err := errors.New(ErrCodeAuth, "Invalid credentials").
		WithUserMessage("Please check your username and password").
		WithContext("attempt", 3).
		WithContext("ip", "192.168.1.1").
		WithWarningSeverity()

	jsonData, _ := json.Marshal(err)

	// Parse to show structure without exact timestamp
	var result map[string]interface{}
	_ = json.Unmarshal(jsonData, &result)

	fmt.Printf("Code: %s\n", result["code"])
	fmt.Printf("Message: %s\n", result["message"])
	fmt.Printf("Severity: %s\n", result["severity"])
	fmt.Printf("User Message: %s\n", result["user_msg"])
	fmt.Printf("Has Context: %t\n", result["context"] != nil)
	fmt.Printf("Has Timestamp: %t\n", result["timestamp"] != nil)
	// Output:
	// Code: AUTHENTICATION_ERROR
	// Message: Invalid credentials
	// Severity: warning
	// User Message: Please check your username and password
	// Has Context: true
	// Has Timestamp: true
}

// ExampleHasCode demonstrates error code checking
func ExampleHasCode() {
	// Create a chain of wrapped errors
	originalErr := fmt.Errorf("disk full")
	wrappedErr := errors.Wrap(originalErr, ErrCodeDatabase, "Failed to write data")
	doubleWrapped := errors.Wrap(wrappedErr, "OPERATION_FAILED", "User save operation failed")

	// Check for specific error codes in the chain
	if errors.HasCode(doubleWrapped, ErrCodeDatabase) {
		fmt.Println("Database error detected in chain")
	}

	if errors.HasCode(doubleWrapped, ErrCodeValidation) {
		fmt.Println("This won't print - no validation error in chain")
	} else {
		fmt.Println("No validation error in chain")
	}
	// Output:
	// Database error detected in chain
	// No validation error in chain
}

// ExampleError_chainMethods demonstrates fluent API usage
func ExampleError_chainMethods() {
	err := errors.New(ErrCodeValidation, "User data validation failed").
		WithUserMessage("Please check your information and try again").
		WithContext("user_id", "12345").
		WithContext("validation_errors", []string{"email", "phone"}).
		AsRetryable().
		WithWarningSeverity()

	fmt.Printf("Code: %s\n", err.ErrorCode())
	fmt.Printf("Retryable: %t\n", err.IsRetryable())
	fmt.Printf("Severity: %s\n", err.Severity)
	// Output:
	// Code: VALIDATION_ERROR
	// Retryable: true
	// Severity: warning
}

// Example of real-world usage patterns
func Example_realWorldUsage() {
	// Simulate a service layer function
	processUser := func(userID string) error {
		// Validate input
		if userID == "" {
			return errors.NewWithField(ErrCodeValidation, "User ID is required", "user_id", userID).
				WithUserMessage("Please provide a valid user ID")
		}

		// Simulate database error
		if userID == "invalid" {
			dbErr := fmt.Errorf("user not found")
			return errors.Wrap(dbErr, ErrCodeDatabase, "Failed to fetch user").
				WithContext("user_id", userID).
				WithUserMessage("User not found")
		}

		return nil
	}

	// Handle errors
	handleError := func(err error) {
		if err == nil {
			return
		}

		structErr, ok := err.(*errors.Error)
		if !ok {
			log.Printf("Unknown error: %v", err)
			return
		}

		// Log technical details
		log.Printf("Error [%s]: %s", structErr.ErrorCode(), structErr.Error())

		// Show user-friendly message
		fmt.Printf("User message: %s\n", structErr.UserMessage())

		// Handle specific error types
		switch {
		case errors.HasCode(err, ErrCodeValidation):
			fmt.Println("Handling validation error")
		case errors.HasCode(err, ErrCodeDatabase):
			fmt.Println("Handling database error")
		}
	}

	// Test cases
	handleError(processUser(""))        // Validation error
	handleError(processUser("invalid")) // Database error
	handleError(processUser("valid"))   // No error

	// Output:
	// User message: Please provide a valid user ID
	// Handling validation error
	// User message: User not found
	// Handling database error
}
