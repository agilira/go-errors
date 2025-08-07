// stacktrace.go: Stacktrace functions for the go-errors AGILira library
//
// Copyright (c) 2025 AGILira
// Series: an AGLIra library
// SPDX-License-Identifier: MPL-2.0

package errors

import (
	"runtime"
	"strconv"
	"strings"
)

// Stacktrace holds a slice of program counters for error tracing.
type Stacktrace struct {
	Frames []uintptr
}

// CaptureStacktrace returns a new Stacktrace from the current call stack.
func CaptureStacktrace(skip int) *Stacktrace {
	const maxDepth = 32
	pcs := make([]uintptr, maxDepth)
	n := runtime.Callers(skip+2, pcs)
	return &Stacktrace{Frames: pcs[:n]}
}

// String returns a human-readable stacktrace.
func (s *Stacktrace) String() string {
	if s == nil || len(s.Frames) == 0 {
		return ""
	}
	var b strings.Builder
	frames := runtime.CallersFrames(s.Frames)
	for {
		frame, more := frames.Next()
		b.WriteString(frame.Function)
		b.WriteString("\n\t")
		b.WriteString(frame.File)
		b.WriteString(":")
		b.WriteString(strconv.Itoa(frame.Line))
		b.WriteString("\n")
		if !more {
			break
		}
	}
	return b.String()
}
