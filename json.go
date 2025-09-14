// json.go: Implementing custom JSON marshaling go-errors AGILira library
//
// Copyright (c) 2025 AGILira - A. Giordano
// Series: an AGLIra library
// SPDX-License-Identifier: MPL-2.0

package errors

import (
	"encoding/json"
)

// MarshalJSON implements custom JSON marshaling for Error.
// It converts the stack trace to a string representation for JSON serialization.
func (e *Error) MarshalJSON() ([]byte, error) {
	type Alias Error
	return json.Marshal(&struct {
		*Alias
		Stack string `json:"stack,omitempty"`
	}{
		Alias: (*Alias)(e),
		Stack: func() string {
			if e.Stack != nil {
				return e.Stack.String()
			}
			return ""
		}(),
	})
}
