// stacktrace.go: Stacktrace functions for the go-errors AGILira library
//
// Copyright (c) 2025 AGILira - A. Giordano
// Series: an AGLIra library
// SPDX-License-Identifier: MPL-2.0

package errors

import (
	"runtime"
	"strconv"
	"strings"
)

// Stacktrace holds a slice of program counters for error tracing and debugging.
// It captures the call stack at the time of error creation for detailed debugging information.
type Stacktrace struct {
	Frames []uintptr
}

// CaptureStacktrace returns a new Stacktrace from the current call stack.
// The skip parameter determines how many stack frames to skip from the top.
// Optimized to reduce allocations by using a smaller initial buffer and growing as needed.
func CaptureStacktrace(skip int) *Stacktrace {
	const (
		initialDepth = 16 // Start smaller - most stacks are shallow
		maxDepth     = 64 // Maximum depth we'll capture
	)

	// Try with initial depth first
	pcs := make([]uintptr, initialDepth)
	n := runtime.Callers(skip+2, pcs)

	// If we filled the buffer, try with larger size
	if n == initialDepth {
		largePcs := make([]uintptr, maxDepth)
		n = runtime.Callers(skip+2, largePcs)
		// Copy only what we need
		result := make([]uintptr, n)
		copy(result, largePcs[:n])
		return &Stacktrace{Frames: result}
	}

	// Copy only the frames we actually captured
	result := make([]uintptr, n)
	copy(result, pcs[:n])
	return &Stacktrace{Frames: result}
}

// String returns a human-readable representation of the stack trace.
// Each frame is displayed with function name, file path, and line number.
// Optimized for better performance with pre-allocated buffer size estimation.
func (s *Stacktrace) String() string {
	if s == nil || len(s.Frames) == 0 {
		return ""
	}

	// Pre-allocate buffer with estimated size to reduce allocations
	// Estimate ~100 chars per frame (function name + file path + line)
	estimatedSize := len(s.Frames) * 100
	var b strings.Builder
	b.Grow(estimatedSize)

	frames := runtime.CallersFrames(s.Frames)
	for {
		frame, more := frames.Next()
		b.WriteString(frame.Function)
		b.WriteString("\n\t")
		b.WriteString(frame.File)
		b.WriteByte(':') // More efficient than WriteString(":")
		b.WriteString(strconv.Itoa(frame.Line))
		b.WriteByte('\n') // More efficient than WriteString("\n")
		if !more {
			break
		}
	}
	return b.String()
}
