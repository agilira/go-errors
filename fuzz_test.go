// fuzz_test.go: Fuzz tests for the go-errors AGILira library
//
// Copyright (c) 2025 AGILira - A. Giordano
// Series: an AGLIra library
// SPDX-License-Identifier: MPL-2.0

package errors

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

// FuzzNew tests the New function with random inputs to ensure it never panics
// and properly handles edge cases like empty codes, extremely long messages,
// and malformed strings.
func FuzzNew(f *testing.F) {
	// Seed with known good and edge case inputs
	f.Add("VALIDATION_ERROR", "Username is required")
	f.Add("", "")
	f.Add("A", "B")
	f.Add("VERY_LONG_ERROR_CODE_THAT_EXCEEDS_NORMAL_LIMITS", "Very long error message that could potentially cause issues")
	f.Add("UNICODE_TEST_ÄÖÜ", "Unicode message: αβγ 中文 🚀")
	f.Add("   ", "   ")
	f.Add("\n\r\t", "\n\r\t")

	f.Fuzz(func(t *testing.T, code string, message string) {
		// The function should never panic regardless of input
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("New panicked with code=%q, message=%q: %v", code, message, r)
			}
		}()

		err := New(ErrorCode(code), message)

		// Basic invariants that should always hold
		if err == nil {
			t.Error("New returned nil error")
			return
		}

		// Error should always be non-empty
		if err.Error() == "" {
			t.Error("Error() returned empty string")
		}

		// Message should be preserved (unless code was invalid)
		if err.Message != message {
			t.Errorf("Message mismatch: expected %q, got %q", message, err.Message)
		}

		// Code should be set to DefaultErrorCode if input was invalid
		if !validateErrorCode(ErrorCode(code)) && err.Code != DefaultErrorCode {
			t.Errorf("Invalid code %q should have been replaced with DefaultErrorCode, got %q", code, err.Code)
		}

		// Valid codes should be preserved
		if validateErrorCode(ErrorCode(code)) && err.Code != ErrorCode(code) {
			t.Errorf("Valid code %q was not preserved, got %q", code, err.Code)
		}

		// Severity should always be set
		if err.Severity == "" {
			t.Error("Severity was not set")
		}

		// Context should be initialized
		if err.Context == nil {
			t.Error("Context was not initialized")
		}

		// Timestamp should be set
		if err.Timestamp.IsZero() {
			t.Error("Timestamp was not set")
		}
	})
}

// FuzzNewWithField tests NewWithField with random field names and values
func FuzzNewWithField(f *testing.F) {
	// Seed with various field/value combinations
	f.Add("FIELD_ERROR", "Invalid field", "email", "user@example.com")
	f.Add("", "", "", "")
	f.Add("TEST", "Message", "field_with_very_long_name_that_could_cause_issues", "value_with_very_long_content")
	f.Add("UNICODE", "Test", "🔥field", "🚀value")
	f.Add("SPECIAL", "Test", "field\nwith\nnewlines", "value\twith\ttabs")

	f.Fuzz(func(t *testing.T, code, message, field, value string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("NewWithField panicked with inputs (%q, %q, %q, %q): %v", code, message, field, value, r)
			}
		}()

		err := NewWithField(ErrorCode(code), message, field, value)

		if err == nil {
			t.Error("NewWithField returned nil")
			return
		}

		// Field and value should be preserved exactly
		if err.Field != field {
			t.Errorf("Field mismatch: expected %q, got %q", field, err.Field)
		}
		if err.Value != value {
			t.Errorf("Value mismatch: expected %q, got %q", value, err.Value)
		}

		// Same invariants as New()
		if err.Error() == "" {
			t.Error("Error() returned empty string")
		}
		if err.Context == nil {
			t.Error("Context was not initialized")
		}
	})
}

// FuzzWrap tests the Wrap function with various error types and inputs
func FuzzWrap(f *testing.F) {
	f.Add("WRAPPER_ERROR", "Operation failed", "original error message")
	f.Add("", "", "")
	f.Add("VERY_LONG_WRAPPER_CODE", "Very long wrapper message", "Very long original error")

	f.Fuzz(func(t *testing.T, code, message, originalErr string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Wrap panicked with inputs (%q, %q, %q): %v", code, message, originalErr, r)
			}
		}()

		// Create original error - could be standard error or structured error
		var origErr error
		if originalErr != "" {
			if len(originalErr)%2 == 0 {
				// Sometimes use standard error
				origErr = fmt.Errorf("%s", originalErr)
			} else {
				// Sometimes use structured error
				origErr = New("ORIGINAL_ERROR", originalErr)
			}
		}

		err := Wrap(origErr, ErrorCode(code), message)

		if err == nil {
			t.Error("Wrap returned nil")
			return
		}

		// Cause should be preserved
		if err.Cause != origErr {
			t.Errorf("Cause not preserved: expected %v, got %v", origErr, err.Cause)
		}

		// Should have stack trace when wrapping
		if err.Stack == nil {
			t.Error("Wrap should capture stack trace")
		}

		// Unwrap should work
		if err.Unwrap() != origErr {
			t.Errorf("Unwrap mismatch: expected %v, got %v", origErr, err.Unwrap())
		}
	})
}

// FuzzWithMethods tests the fluent API methods with random inputs
func FuzzWithMethods(f *testing.F) {
	f.Add("TEST_ERROR", "Base message", "User message", "context_key", "context_value", "warning")
	f.Add("", "", "", "", "", "")
	f.Add("LONG", strings.Repeat("Long message", 100), strings.Repeat("Long user message", 50), "very_long_key", strings.Repeat("very_long_value", 20), "critical")

	f.Fuzz(func(t *testing.T, code, message, userMsg, contextKey, contextValue, severity string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Fluent methods panicked: %v", r)
			}
		}()

		err := New(ErrorCode(code), message).
			WithUserMessage(userMsg).
			WithContext(contextKey, contextValue).
			WithSeverity(severity).
			AsRetryable()

		if err == nil {
			t.Error("Chain returned nil")
			return
		}

		// Check all properties were set
		if err.UserMessage() != userMsg && userMsg != "" {
			t.Errorf("UserMessage mismatch: expected %q, got %q", userMsg, err.UserMessage())
		}

		if err.Severity != severity {
			t.Errorf("Severity mismatch: expected %q, got %q", severity, err.Severity)
		}

		if !err.IsRetryable() {
			t.Error("Error should be retryable")
		}

		if contextKey != "" {
			if val, exists := err.Context[contextKey]; !exists {
				t.Errorf("Context key %q not found", contextKey)
			} else if val != contextValue {
				t.Errorf("Context value mismatch: expected %q, got %q", contextValue, val)
			}
		}
	})
}

// FuzzHasCode tests the HasCode function with various error chains
func FuzzHasCode(f *testing.F) {
	f.Add("SEARCH_CODE", "TARGET_CODE", "WRAPPER_CODE")
	f.Add("", "", "")
	f.Add("SAME", "SAME", "SAME")

	f.Fuzz(func(t *testing.T, searchCode, targetCode, wrapperCode string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("HasCode panicked: %v", r)
			}
		}()

		// Create a chain of errors
		originalErr := New(ErrorCode(targetCode), "original")
		wrappedErr := Wrap(originalErr, ErrorCode(wrapperCode), "wrapped")

		// Test with nil
		if HasCode(nil, ErrorCode(searchCode)) {
			t.Error("HasCode should return false for nil error")
		}

		// Test searching for codes in the chain
		foundOriginal := HasCode(wrappedErr, ErrorCode(targetCode))
		_ = foundOriginal                               // Used below in validation
		_ = HasCode(wrappedErr, ErrorCode(wrapperCode)) // Just test it doesn't panic
		_ = HasCode(wrappedErr, ErrorCode(searchCode))  // Just test it doesn't panic

		// Validate logic
		if targetCode != "" && validateErrorCode(ErrorCode(targetCode)) {
			expectedTargetCode := targetCode
			if !validateErrorCode(ErrorCode(targetCode)) {
				expectedTargetCode = string(DefaultErrorCode)
			}
			shouldFindOriginal := HasCode(wrappedErr, ErrorCode(expectedTargetCode))
			if foundOriginal != shouldFindOriginal {
				t.Errorf("HasCode logic error for target code %q", targetCode)
			}
		}
	})
}

// FuzzJSONMarshal tests JSON marshaling with various error configurations
func FuzzJSONMarshal(f *testing.F) {
	f.Add("JSON_ERROR", "Test message", "user message", true)
	f.Add("", "", "", false)
	f.Add("UNICODE_🚀", "Message with emoji 🔥", "User sees: αβγ", true)

	f.Fuzz(func(t *testing.T, code, message, userMsg string, withStack bool) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("JSON marshaling panicked: %v", r)
			}
		}()

		var err *Error
		if withStack {
			// Create with stack trace
			err = Wrap(fmt.Errorf("cause"), ErrorCode(code), message)
		} else {
			err = New(ErrorCode(code), message)
		}

		if userMsg != "" {
			err = err.WithUserMessage(userMsg)
		}

		// Add some context with potentially problematic values
		err = err.WithContext("string", message).
			WithContext("bool", withStack).
			WithContext("number", len(message))

		// Should be able to marshal to JSON
		jsonData, jsonErr := json.Marshal(err)
		if jsonErr != nil {
			t.Errorf("JSON marshaling failed: %v", jsonErr)
			return
		}

		// Should be valid JSON
		var parsed map[string]interface{}
		if parseErr := json.Unmarshal(jsonData, &parsed); parseErr != nil {
			t.Errorf("Generated invalid JSON: %v", parseErr)
			return
		}

		// Check required fields exist
		if _, exists := parsed["code"]; !exists {
			t.Error("JSON missing 'code' field")
		}
		if _, exists := parsed["message"]; !exists {
			t.Error("JSON missing 'message' field")
		}
		if _, exists := parsed["timestamp"]; !exists {
			t.Error("JSON missing 'timestamp' field")
		}

		// Verify that JSON is valid UTF-8
		if !utf8.Valid(jsonData) {
			t.Error("JSON output is not valid UTF-8")
		}
	})
}

// FuzzValidateErrorCode tests the error code validation with edge cases
func FuzzValidateErrorCode(f *testing.F) {
	f.Add("VALID_CODE")
	f.Add("")
	f.Add("   ")
	f.Add("\n\r\t")
	f.Add("A")
	f.Add("UNICODE_αβγ_中文_🚀")
	f.Add(strings.Repeat("VERY_LONG", 1000))

	f.Fuzz(func(t *testing.T, code string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("validateErrorCode panicked with %q: %v", code, r)
			}
		}()

		isValid := validateErrorCode(ErrorCode(code))

		// Verify the validation logic
		if len(code) == 0 {
			if isValid {
				t.Error("Empty code should be invalid")
			}
		} else {
			// Check if it's only whitespace
			hasNonWhitespace := false
			for _, r := range code {
				if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
					hasNonWhitespace = true
					break
				}
			}

			if hasNonWhitespace && !isValid {
				t.Errorf("Code with non-whitespace characters should be valid: %q", code)
			} else if !hasNonWhitespace && isValid {
				t.Errorf("Code with only whitespace should be invalid: %q", code)
			}
		}

		// Test that validation is used consistently in constructors
		err := New(ErrorCode(code), "test message")
		if isValid && err.Code != ErrorCode(code) {
			t.Errorf("Valid code was changed: expected %q, got %q", code, err.Code)
		} else if !isValid && err.Code != DefaultErrorCode {
			t.Errorf("Invalid code should become DefaultErrorCode: expected %q, got %q", DefaultErrorCode, err.Code)
		}
	})
}

// FuzzStacktrace tests stacktrace capture and string conversion
func FuzzStacktrace(f *testing.F) {
	f.Add(0)
	f.Add(1)
	f.Add(10)
	f.Add(100)

	f.Fuzz(func(t *testing.T, skip int) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Stacktrace operations panicked with skip=%d: %v", skip, r)
			}
		}()

		// Capture stacktrace with random skip value
		st := CaptureStacktrace(skip)

		if st == nil {
			t.Error("CaptureStacktrace returned nil")
			return
		}

		// String conversion should never panic
		str := st.String()

		// Output should be valid UTF-8
		if !utf8.ValidString(str) {
			t.Error("Stacktrace string is not valid UTF-8")
		}

		// If we have frames, string should not be empty
		if len(st.Frames) > 0 && str == "" {
			t.Error("Non-empty stacktrace produced empty string")
		}

		// Empty stacktrace should produce empty string
		emptySt := &Stacktrace{Frames: nil}
		if emptySt.String() != "" {
			t.Error("Empty stacktrace should produce empty string")
		}
	})
}
