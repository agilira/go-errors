// errors_test.go: Tests for the go-errors AGILira library
//
// Copyright (c) 2025 AGILira
// Series: an AGLIra library
// SPDX-License-Identifier: MPL-2.0

package errors

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

const (
	TestCodeValidation ErrorCode = "VALIDATION_ERROR"
	TestCodeDatabase   ErrorCode = "DATABASE_ERROR"
)

func TestNew(t *testing.T) {
	err := New(TestCodeValidation, "Validation failed")
	if err.Code != TestCodeValidation {
		t.Errorf("Expected code %s, got %s", TestCodeValidation, err.Code)
	}
	if err.Message != "Validation failed" {
		t.Errorf("Expected message 'Validation failed', got '%s'", err.Message)
	}
	if err.Severity != "error" {
		t.Errorf("Expected severity 'error', got '%s'", err.Severity)
	}
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

	data, errMarshal := json.Marshal(err)
	if errMarshal != nil {
		t.Errorf("MarshalJSON failed: %v", errMarshal)
	}

	// Verify that stack field is empty string when nil
	var result map[string]interface{}
	if errUnmarshal := json.Unmarshal(data, &result); errUnmarshal != nil {
		t.Errorf("Failed to unmarshal: %v", errUnmarshal)
	}

	// Stack should be empty string when nil
	if stack, exists := result["stack"]; exists && stack != "" {
		t.Error("Stack should be empty string when nil")
	}
}

func TestMarshalJSONWithEmptyStack(t *testing.T) {
	// Create error with empty stack trace
	err := New(TestCodeValidation, "Validation failed")
	err.Stack = &Stacktrace{Frames: []uintptr{}} // Empty stack

	data, errMarshal := json.Marshal(err)
	if errMarshal != nil {
		t.Errorf("MarshalJSON failed: %v", errMarshal)
	}

	// This should force the execution of the return "" branch
	var result map[string]interface{}
	if errUnmarshal := json.Unmarshal(data, &result); errUnmarshal != nil {
		t.Errorf("Failed to unmarshal: %v", errUnmarshal)
	}

	// Stack should be empty string when stack has no frames
	if stack, exists := result["stack"]; exists && stack != "" {
		t.Error("Stack should be empty string when stack has no frames")
	}
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
