// errors_test.go: Tests for the go-errors AGILira library
//
// Copyright (c) 2025 AGILira - A. Giordano
// Series: an AGLIra library
// SPDX-License-Identifier: MPL-2.0

package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

const (
	TestCodeValidation ErrorCode = "VALIDATION_ERROR"
	TestCodeDatabase   ErrorCode = "DATABASE_ERROR"
)

func TestNew(t *testing.T) {
	err := New(TestCodeValidation, "Validation failed")

	// Use a table-driven approach to reduce cyclomatic complexity
	tests := []struct {
		field    string
		expected interface{}
		actual   interface{}
	}{
		{"Code", TestCodeValidation, err.Code},
		{"Message", "Validation failed", err.Message},
		{"Severity", SeverityError, err.Severity},
	}

	for _, tt := range tests {
		if tt.actual != tt.expected {
			t.Errorf("Field %s: expected %v, got %v", tt.field, tt.expected, tt.actual)
		}
	}

	// Test special cases separately for clarity
	if err.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}
	if err.Context == nil {
		t.Error("Expected context to be initialized")
	}
}

func TestNewWithField(t *testing.T) {
	err := NewWithField(TestCodeValidation, "Field is empty", "username", "")
	if err.Field != "username" {
		t.Errorf("Expected field 'username', got '%s'", err.Field)
	}
	if err.Value != "" {
		t.Errorf("Expected empty value, got '%s'", err.Value)
	}
}

func TestNewWithContext(t *testing.T) {
	ctx := map[string]interface{}{"user_id": "123"}
	err := NewWithContext(TestCodeDatabase, "DB error", ctx)
	if err.Context["user_id"] != "123" {
		t.Errorf("Expected user_id '123', got '%v'", err.Context["user_id"])
	}
}

func TestWrapAndRootCause(t *testing.T) {
	orig := errors.New("original error")
	wrapped := Wrap(orig, TestCodeDatabase, "DB failed")
	if wrapped.Cause != orig {
		t.Error("Expected cause to be original error")
	}
	if RootCause(wrapped) != orig {
		t.Error("RootCause did not return the original error")
	}
}

func TestErrorString(t *testing.T) {
	err := NewWithField(TestCodeValidation, "Invalid format", "email", "bad")
	str := err.Error()
	if !strings.Contains(str, string(TestCodeValidation)) || !strings.Contains(str, "Invalid format") {
		t.Error("Error string does not contain code or message")
	}
}

func TestUnwrap(t *testing.T) {
	orig := errors.New("original error")
	wrapped := Wrap(orig, TestCodeDatabase, "DB failed")
	if wrapped.Unwrap() != orig {
		t.Error("Expected unwrapped error to be original error")
	}
}

func TestHasCode(t *testing.T) {
	err := Wrap(New(TestCodeValidation, "Validation failed"), TestCodeDatabase, "DB failed")
	if !HasCode(err, TestCodeDatabase) {
		t.Error("HasCode did not find the code in the chain")
	}
	if HasCode(err, "NON_EXISTENT") {
		t.Error("HasCode returned true for a non-existent code")
	}
}

func TestIsAndAsCompatibility(t *testing.T) {
	err := New(TestCodeValidation, "Validation failed")
	if !errors.Is(err, &Error{Code: TestCodeValidation}) {
		t.Error("errors.Is did not match error code")
	}
	var target *Error
	if !errors.As(err, &target) {
		t.Error("errors.As did not match *Error type")
	}
}

func TestStacktrace(t *testing.T) {
	err := Wrap(errors.New("fail"), TestCodeValidation, "with stack")
	if err.Stack == nil || len(err.Stack.Frames) == 0 {
		t.Error("Expected stacktrace to be captured")
	}
	str := err.Stack.String()
	if str == "" {
		t.Error("Stacktrace string should not be empty")
	}
}

func TestMarshalJSON(t *testing.T) {
	err := New(TestCodeValidation, "Validation failed")
	data, errMarshal := json.Marshal(err)
	if errMarshal != nil {
		t.Errorf("MarshalJSON failed: %v", errMarshal)
	}
	if !strings.Contains(string(data), "Validation failed") {
		t.Error("JSON does not contain error message")
	}
}

func TestInterfaces(t *testing.T) {
	err := New(TestCodeValidation, "Validation failed")

	// Test ErrorCoder interface
	if err.ErrorCode() != TestCodeValidation {
		t.Error("ErrorCoder interface not working")
	}

	err.Retryable = true
	// Test Retryable interface
	if !err.IsRetryable() {
		t.Error("Retryable interface not working")
	}

	err = err.WithUserMessage("User message")
	// Test UserMessager interface
	if err.UserMessage() != "User message" {
		t.Error("UserMessager interface not working")
	}
}

func TestIsWithNilTarget(t *testing.T) {
	err := New(TestCodeValidation, "Validation failed")
	if err.Is(nil) {
		t.Error("Is should return false for nil target")
	}
}

func TestIsWithNonErrorTarget(t *testing.T) {
	err := New(TestCodeValidation, "Validation failed")
	nonErrorTarget := errors.New("standard error")
	if err.Is(nonErrorTarget) {
		t.Error("Is should return false for non-Error target")
	}
}

func TestAsFunction(t *testing.T) {
	orig := errors.New("original error")
	wrapped := Wrap(orig, TestCodeDatabase, "DB failed")

	// Test that As can find the original error in the cause chain
	var origTarget error
	if !wrapped.As(&origTarget) {
		t.Error("As should succeed for error target")
	}
	if origTarget != orig {
		t.Error("As should set target to the original error")
	}
}

func TestAsWithNilTarget(t *testing.T) {
	err := New(TestCodeValidation, "Validation failed")
	if err.As(nil) {
		t.Error("As should return false for nil target")
	}
}

func TestMarshalJSONWithNilStack(t *testing.T) {
	// Create error without stack trace (should be nil by default)
	err := New(TestCodeValidation, "Validation failed")

	// Verify stack is actually nil
	if err.Stack != nil {
		t.Fatal("Stack should be nil for errors created with New()")
	}

	// Test JSON marshaling and validation in helper function to reduce complexity
	validateJSONStack(t, err, "Stack should be empty string when nil")
}

// Helper function to reduce cyclomatic complexity in JSON tests
func validateJSONStack(t *testing.T, err *Error, expectedMsg string) {
	data, errMarshal := json.Marshal(err)
	if errMarshal != nil {
		t.Errorf("MarshalJSON failed: %v", errMarshal)
		return
	}

	var result map[string]interface{}
	if errUnmarshal := json.Unmarshal(data, &result); errUnmarshal != nil {
		t.Errorf("Failed to unmarshal: %v", errUnmarshal)
		return
	}

	if stack, exists := result["stack"]; exists && stack != "" {
		t.Error(expectedMsg)
	}
}

func TestMarshalJSONWithEmptyStack(t *testing.T) {
	// Create error with empty stack trace
	err := New(TestCodeValidation, "Validation failed")
	err.Stack = &Stacktrace{Frames: []uintptr{}} // Empty stack

	// Use helper function to reduce complexity
	validateJSONStack(t, err, "Stack should be empty string when stack has no frames")
}

func TestStacktraceStringWithNilStacktrace(t *testing.T) {
	var stack *Stacktrace
	result := stack.String()
	if result != "" {
		t.Errorf("Expected empty string for nil stacktrace, got: %s", result)
	}
}

func TestStacktraceStringWithEmptyFrames(t *testing.T) {
	stack := &Stacktrace{Frames: []uintptr{}}
	result := stack.String()
	if result != "" {
		t.Errorf("Expected empty string for empty frames, got: %s", result)
	}
}

func TestUserMessageWithEmptyUserMsg(t *testing.T) {
	err := New(TestCodeValidation, "Technical message")
	err.UserMsg = "" // Ensure UserMsg is empty

	result := err.UserMessage()
	if result != "Technical message" {
		t.Errorf("Expected technical message when UserMsg is empty, got: %s", result)
	}
}

func TestUserMessageWithEmptyBoth(t *testing.T) {
	err := New(TestCodeValidation, "")
	err.UserMsg = ""

	result := err.UserMessage()
	if result != "" {
		t.Errorf("Expected empty string when both messages are empty, got: %s", result)
	}
}

func TestWithContext(t *testing.T) {
	err := New(TestCodeValidation, "Test error")

	// Test adding context to error with nil context
	err.Context = nil
	err = err.WithContext("key1", "value1")
	if err.Context == nil {
		t.Error("Expected context to be initialized")
	}
	if err.Context["key1"] != "value1" {
		t.Errorf("Expected context key1 to be 'value1', got '%v'", err.Context["key1"])
	}

	// Test adding more context
	err = err.WithContext("key2", "value2")
	if err.Context["key2"] != "value2" {
		t.Errorf("Expected context key2 to be 'value2', got '%v'", err.Context["key2"])
	}

	// Test updating existing context
	err = err.WithContext("key1", "updated_value")
	if err.Context["key1"] != "updated_value" {
		t.Errorf("Expected context key1 to be 'updated_value', got '%v'", err.Context["key1"])
	}
}

func TestAsRetryable(t *testing.T) {
	err := New(TestCodeValidation, "Test error")

	// Test marking error as retryable
	err = err.AsRetryable()
	if !err.Retryable {
		t.Error("Expected error to be marked as retryable")
	}

	// Test that it returns the same error instance for chaining
	if err.AsRetryable() != err {
		t.Error("Expected AsRetryable to return the same error instance")
	}
}

func TestWithSeverity(t *testing.T) {
	err := New(TestCodeValidation, "Test error")

	// Test setting severity
	err = err.WithSeverity("warning")
	if err.Severity != "warning" {
		t.Errorf("Expected severity to be 'warning', got '%s'", err.Severity)
	}

	// Test setting different severity
	err = err.WithSeverity("critical")
	if err.Severity != "critical" {
		t.Errorf("Expected severity to be 'critical', got '%s'", err.Severity)
	}

	// Test that it returns the same error instance for chaining
	if err.WithSeverity("info") != err {
		t.Error("Expected WithSeverity to return the same error instance")
	}
}

func TestMethodChaining(t *testing.T) {
	err := New(TestCodeValidation, "Test").
		WithUserMessage("User friendly message").
		WithContext("key", "value").
		AsRetryable().
		WithSeverity("warning")

	if err.UserMsg != "User friendly message" {
		t.Error("Method chaining failed for UserMessage")
	}
	if err.Context["key"] != "value" {
		t.Error("Method chaining failed for Context")
	}
	if !err.Retryable {
		t.Error("Method chaining failed for Retryable")
	}
	if err.Severity != "warning" {
		t.Error("Method chaining failed for Severity")
	}
}

func TestSeverityConstants(t *testing.T) {
	// Test that severity constants are properly defined
	if SeverityCritical != "critical" {
		t.Errorf("Expected SeverityCritical to be 'critical', got '%s'", SeverityCritical)
	}
	if SeverityError != "error" {
		t.Errorf("Expected SeverityError to be 'error', got '%s'", SeverityError)
	}
	if SeverityWarning != "warning" {
		t.Errorf("Expected SeverityWarning to be 'warning', got '%s'", SeverityWarning)
	}
	if SeverityInfo != "info" {
		t.Errorf("Expected SeverityInfo to be 'info', got '%s'", SeverityInfo)
	}
}

func TestSeverityHelperMethods(t *testing.T) {
	err := New(TestCodeValidation, "Test message")

	// Test critical severity helper (return value intentionally ignored to test in-place modification)
	_ = err.WithCriticalSeverity()
	if err.Severity != SeverityCritical {
		t.Errorf("Expected severity '%s', got '%s'", SeverityCritical, err.Severity)
	}

	// Test warning severity helper (return value intentionally ignored to test in-place modification)
	_ = err.WithWarningSeverity()
	if err.Severity != SeverityWarning {
		t.Errorf("Expected severity '%s', got '%s'", SeverityWarning, err.Severity)
	}

	// Test info severity helper (return value intentionally ignored to test in-place modification)
	_ = err.WithInfoSeverity()
	if err.Severity != SeverityInfo {
		t.Errorf("Expected severity '%s', got '%s'", SeverityInfo, err.Severity)
	}
}

func TestErrorCodeValidation(t *testing.T) {
	// Table-driven tests to reduce cyclomatic complexity
	testCases := []struct {
		name         string
		input        ErrorCode
		expectedCode ErrorCode
	}{
		{"empty ErrorCode", "", DefaultErrorCode},
		{"whitespace-only ErrorCode", "   \t\n  ", DefaultErrorCode},
		{"valid ErrorCode", "VALID_CODE", "VALID_CODE"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := New(tc.input, "Test message")
			if err.Code != tc.expectedCode {
				t.Errorf("Expected code '%s', got '%s'", tc.expectedCode, err.Code)
			}
		})
	}

	// Test constructor functions with empty codes
	testConstructorValidation(t)
}

// Helper function to test constructor validation
func testConstructorValidation(t *testing.T) {
	constructors := []struct {
		name string
		fn   func() ErrorCode
	}{
		{"NewWithField", func() ErrorCode {
			return NewWithField("", "Test message", "field", "value").Code
		}},
		{"NewWithContext", func() ErrorCode {
			ctx := map[string]interface{}{"key": "value"}
			return NewWithContext("", "Test message", ctx).Code
		}},
		{"Wrap", func() ErrorCode {
			originalErr := fmt.Errorf("original error")
			return Wrap(originalErr, "", "Wrapped message").Code
		}},
	}

	for _, constructor := range constructors {
		t.Run(constructor.name, func(t *testing.T) {
			code := constructor.fn()
			if code != DefaultErrorCode {
				t.Errorf("Expected code '%s' for empty input in %s, got '%s'",
					DefaultErrorCode, constructor.name, code)
			}
		})
	}
}

func TestStacktraceOptimizations(t *testing.T) {
	// Test stacktrace capture with normal depth
	stack := CaptureStacktrace(0)
	if stack == nil {
		t.Fatal("Expected non-nil stacktrace")
	}

	if len(stack.Frames) == 0 {
		t.Error("Expected frames in stacktrace")
	}

	// Test String() method optimizations
	stackStr := stack.String()
	if len(stackStr) == 0 {
		t.Error("Expected non-empty stacktrace string")
	}

	// Test direct access to optimization paths by calling recursive function
	recursiveStack := testRecursiveStackCapture(0, 20) // Create deeper stack
	if recursiveStack == nil {
		t.Fatal("Expected non-nil recursive stacktrace")
	}
}

// Helper function to create deeper call stacks for testing
func testRecursiveStackCapture(depth, target int) *Stacktrace {
	if depth >= target {
		return CaptureStacktrace(0)
	}
	return testRecursiveStackCapture(depth+1, target)
}
