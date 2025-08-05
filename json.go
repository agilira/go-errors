// json.go: Iplementing custom JSON marshaling go-errors AGILira library
//
// Copyright (c) 2025 AGILira
// Series: an AGLIra fragment
// SPDX-License-Identifier: MPL-2.0

package errors

import (
	"encoding/json"
)

// MarshalJSON implements custom JSON marshaling for Error.
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
