// benchmark_test.go: Benchmarks for the go-errors AGILira library
//
// Copyright (c) 2025 AGILira - A. Giordano
// Series: an AGLIra library
// SPDX-License-Identifier: MPL-2.0

package errors

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/agilira/go-timecache"
)

const BenchmarkErrorCode ErrorCode = "BENCHMARK_ERROR"

// Benchmark error creation functions
func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = New(BenchmarkErrorCode, "Benchmark error message")
	}
}

func BenchmarkNewWithField(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewWithField(BenchmarkErrorCode, "Validation error", "username", "invalid")
	}
}

func BenchmarkNewWithContext(b *testing.B) {
	ctx := map[string]interface{}{
		"user_id": "123",
		"action":  "login",
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewWithContext(BenchmarkErrorCode, "Context error", ctx)
	}
}

func BenchmarkWrap(b *testing.B) {
	originalErr := fmt.Errorf("original error")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Wrap(originalErr, BenchmarkErrorCode, "Wrapped error")
	}
}

// Benchmark stacktrace operations
func BenchmarkCaptureStacktrace(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = CaptureStacktrace(1)
	}
}

func BenchmarkStacktraceString(b *testing.B) {
	stack := CaptureStacktrace(1)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = stack.String()
	}
}

// Benchmark JSON marshaling
func BenchmarkMarshalJSON(b *testing.B) {
	err := New(BenchmarkErrorCode, "JSON benchmark error").
		WithUserMessage("User friendly message").
		WithContext("key", "value").
		AsRetryable().
		WithCriticalSeverity()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(err)
	}
}

func BenchmarkMarshalJSONWithStack(b *testing.B) {
	originalErr := fmt.Errorf("original error")
	err := Wrap(originalErr, BenchmarkErrorCode, "Error with stack").
		WithUserMessage("User friendly message").
		WithContext("operation", "benchmark")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(err)
	}
}

// Benchmark method chaining
func BenchmarkMethodChaining(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = New(BenchmarkErrorCode, "Chaining benchmark").
			WithUserMessage("User message").
			WithContext("iteration", i).
			WithContext("benchmark", true).
			AsRetryable().
			WithWarningSeverity()
	}
}

// Benchmark error interface methods
func BenchmarkErrorString(b *testing.B) {
	err := New(BenchmarkErrorCode, "Error string benchmark")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkHasCode(b *testing.B) {
	err := New(BenchmarkErrorCode, "HasCode benchmark")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = HasCode(err, BenchmarkErrorCode)
	}
}

func BenchmarkRootCause(b *testing.B) {
	originalErr := fmt.Errorf("original error")
	wrappedErr := Wrap(originalErr, BenchmarkErrorCode, "Wrapped error")
	doubleWrapped := Wrap(wrappedErr, "DOUBLE_WRAP", "Double wrapped")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = RootCause(doubleWrapped)
	}
}

// Benchmark validation functions
func BenchmarkValidateErrorCode(b *testing.B) {
	testCodes := []ErrorCode{
		"VALID_CODE",
		"",
		"   \t\n  ",
		"ANOTHER_VALID_CODE",
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, code := range testCodes {
			_ = validateErrorCode(code)
		}
	}
}

// Benchmark time access methods comparison
func BenchmarkTimeNow(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = time.Now()
	}
}

func BenchmarkTimeCachedTime(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = timecache.CachedTime()
	}
}
