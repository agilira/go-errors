# API Reference

This document provides a complete reference for all public APIs in the go-errors library.

## Package Overview

The `errors` package provides structured, contextual error handling for Go applications. It offers a comprehensive error handling system with features like structured error types, stack traces, user-friendly messages, JSON serialization, retry logic, and interface-based error handling.

## Types

### ErrorCode
```go
type ErrorCode string
```
ErrorCode represents a custom error code for categorization and programmatic error handling.

### Error
```go
type Error struct {
    Code      ErrorCode              `json:"code"`
    Message   string                 `json:"message"`
    Field     string                 `json:"field,omitempty"`
    Value     string                 `json:"value,omitempty"`
    Context   map[string]interface{} `json:"context,omitempty"`
    Timestamp time.Time              `json:"timestamp"`
    Cause     error                  `json:"cause,omitempty"`
    Severity  string                 `json:"severity"`
    Stack     *Stacktrace            `json:"stack,omitempty"`
    UserMsg   string                 `json:"user_msg,omitempty"`
    Retryable bool                   `json:"retryable,omitempty"`
}
```

**Fields:**
- `Code`: Application-defined error code for categorization
- `Message`: Technical error message for developers
- `Field`: Field name for validation errors
- `Value`: Field value for validation errors
- `Context`: Additional structured data
- `Timestamp`: Error creation time
- `Cause`: Underlying error in the chain
- `Severity`: Error severity level (default: "error")
- `Stack`: Captured stack trace (if available)
- `UserMsg`: User-friendly error message
- `Retryable`: Indicates if the error can be retried

### Stacktrace
```go
type Stacktrace struct {
    Frames []uintptr
}
```
Stacktrace holds program counters for error tracing and debugging.

## Constructor Functions

### New
```go
func New(code ErrorCode, message string) *Error
```
Creates a new error with the specified code and message.

**Parameters:**
- `code`: Error code for categorization
- `message`: Technical error message

**Returns:** Pointer to a new Error instance

**Example:**
```go
err := errors.New("VALIDATION_ERROR", "Invalid input format")
```

### NewWithField
```go
func NewWithField(code ErrorCode, message, field, value string) *Error
```
Creates a new error with field validation context.

**Parameters:**
- `code`: Error code for categorization
- `message`: Technical error message
- `field`: Field name that failed validation
- `value`: Field value that failed validation

**Returns:** Pointer to a new Error instance

**Example:**
```go
err := errors.NewWithField("VALIDATION_ERROR", "Email format invalid", "email", "invalid-email")
```

### NewWithContext
```go
func NewWithContext(code ErrorCode, message string, context map[string]interface{}) *Error
```
Creates a new error with additional structured context.

**Parameters:**
- `code`: Error code for categorization
- `message`: Technical error message
- `context`: Additional structured data

**Returns:** Pointer to a new Error instance

**Example:**
```go
ctx := map[string]interface{}{
    "user_id": "123",
    "operation": "create",
}
err := errors.NewWithContext("DATABASE_ERROR", "Insert failed", ctx)
```

## Helper Functions

These functions provide utility operations for error handling and manipulation.

### Wrap
```go
func Wrap(err error, code ErrorCode, message string) *Error
```
Wraps an existing error with additional context and captures stack trace.

**Parameters:**
- `err`: Original error to wrap
- `code`: Error code for the wrapper
- `message`: Additional context message

**Returns:** Pointer to a new Error instance

**Example:**
```go
originalErr := errors.New("DB_ERROR", "Connection failed")
wrappedErr := errors.Wrap(originalErr, "SERVICE_ERROR", "User creation failed")
```

### RootCause
```go
func RootCause(err error) error
```
Returns the original error in the error chain.

**Parameters:**
- `err`: Error to analyze

**Returns:** The root cause error

**Example:**
```go
root := errors.RootCause(wrappedErr)
```

### HasCode
```go
func HasCode(err error, code ErrorCode) bool
```
Checks if any error in the chain has the specified code.

**Parameters:**
- `err`: Error to check
- `code`: Error code to search for

**Returns:** True if the code is found in the error chain

**Example:**
```go
if errors.HasCode(err, "VALIDATION_ERROR") {
    // Handle validation error
}
```

## Error Methods

These methods are available on *Error instances and implement standard interfaces.

### Error
```go
func (e *Error) Error() string
```
Implements the error interface. Returns formatted error string.

**Returns:** Formatted error message in format "[CODE]: message"

### Unwrap
```go
func (e *Error) Unwrap() error
```
Implements the error unwrapping interface.

**Returns:** The underlying cause error

### Is
```go
func (e *Error) Is(target error) bool
```
Implements errors.Is compatibility. Compares error codes between Error instances.

**Parameters:**
- `target`: Target error to compare against

**Returns:** True if both errors are *Error instances with matching codes

**Example:**
```go
err1 := errors.New("VALIDATION_ERROR", "Invalid input")
err2 := errors.New("VALIDATION_ERROR", "Different message")
if err1.Is(err2) {
    // This will be true - same error code
}
```

### As
```go
func (e *Error) As(target interface{}) bool
```
Implements errors.As compatibility. Delegates to the standard library's errors.As on the cause.

**Parameters:**
- `target`: Target interface to search for

**Returns:** True if target is found in the cause chain using standard library logic

**Note:** This method searches the cause chain, not the error itself

### WithUserMessage
```go
func (e *Error) WithUserMessage(msg string) *Error
```
Sets a user-friendly message on the error and returns the error for chaining.

**Parameters:**
- `msg`: User-friendly error message

**Returns:** Self-reference for method chaining

**Example:**
```go
err := errors.New("VALIDATION_ERROR", "Invalid input").WithUserMessage("Please check your input and try again")
```

### WithContext
```go
func (e *Error) WithContext(key string, value interface{}) *Error
```
Adds or updates context information on the error and returns the error for chaining.

**Parameters:**
- `key`: Context key
- `value`: Context value

**Returns:** Self-reference for method chaining

**Example:**
```go
err := errors.New("DATABASE_ERROR", "Query failed").WithContext("query", "SELECT * FROM users")
```

### AsRetryable
```go
func (e *Error) AsRetryable() *Error
```
Marks the error as retryable and returns the error for chaining.

**Returns:** Self-reference for method chaining

**Example:**
```go
err := errors.New("NETWORK_ERROR", "Connection timeout").AsRetryable()
```

### WithSeverity
```go
func (e *Error) WithSeverity(severity string) *Error
```
Sets the severity level of the error and returns the error for chaining.

**Parameters:**
- `severity`: Severity level ("error", "warning", "info", "critical")

**Returns:** Self-reference for method chaining

**Example:**
```go
err := errors.New("VALIDATION_ERROR", "Invalid input").WithSeverity("warning")
```

### UserMessage
```go
func (e *Error) UserMessage() string
```
Returns the user-friendly message if set, otherwise falls back to the technical message.

**Returns:** User-friendly or technical message

### ErrorCode
```go
func (e *Error) ErrorCode() ErrorCode
```
Returns the error code.

**Returns:** The error code

### IsRetryable
```go
func (e *Error) IsRetryable() bool
```
Returns whether the error is retryable.

**Returns:** True if the error is retryable

### MarshalJSON
```go
func (e *Error) MarshalJSON() ([]byte, error)
```
Implements custom JSON marshaling for Error.

**Returns:** JSON bytes and error

## Stacktrace Methods

### CaptureStacktrace
```go
func CaptureStacktrace(skip int) *Stacktrace
```
Captures the current call stack.

**Parameters:**
- `skip`: Number of frames to skip (typically 1 for immediate caller)

**Returns:** Pointer to Stacktrace instance

**Implementation Details:**
- Captures up to 32 frames maximum
- Uses runtime.Callers for efficient stack capture
- Automatically skips 2 additional frames for internal processing

### String
```go
func (s *Stacktrace) String() string
```
Returns a human-readable stack trace.

**Returns:** Formatted stack trace string or empty string if stack is nil or empty

**Format:** Each frame includes function name, file path, and line number

## Interfaces

### ErrorCoder
```go
type ErrorCoder interface {
    ErrorCode() ErrorCode
}
```
Interface for extracting error codes from errors.

**Implementation:** The Error type implements this interface by returning its Code field.

### Retryable
```go
type Retryable interface {
    IsRetryable() bool
}
```
Interface for checking if an error is retryable.

**Implementation:** The Error type implements this interface by returning its Retryable field.

### UserMessager
```go
type UserMessager interface {
    UserMessage() string
}
```
Interface for extracting user-friendly messages from errors.

**Implementation:** The Error type implements this interface by returning the user message or falling back to the technical message.

## Interface Implementations

The Error type implements all three interfaces:

```go
// ErrorCoder implementation
func (e *Error) ErrorCode() ErrorCode {
    return e.Code
}

// Retryable implementation  
func (e *Error) IsRetryable() bool {
    return e.Retryable
}

// UserMessager implementation
func (e *Error) UserMessage() string {
    if e.UserMsg != "" {
        return e.UserMsg
    }
    return e.Message
}
```

**Note:** These methods are implemented in the `usermsg.go` file alongside the user message functionality.

## Error Handling Patterns

### Standard Library Compatibility
The library is fully compatible with Go's standard error handling:

```go
// errors.Is compatibility
if errors.Is(err, targetErr) {
    // Handle specific error
}

// errors.As compatibility
var apiErr *errors.Error
if errors.As(err, &apiErr) {
    // Handle structured error
}

// errors.Unwrap compatibility
cause := errors.Unwrap(err)
```

**Important Notes:**
- `Is` method compares error codes between *Error instances only
- `As` method delegates to standard library's errors.As on the cause chain
- `Unwrap` returns the immediate cause, not the root cause

### Error Code Constants
Define error codes as constants for consistency:

```go
const (
    ErrCodeValidation = "VALIDATION_ERROR"
    ErrCodeDatabase   = "DATABASE_ERROR"
    ErrCodeNetwork    = "NETWORK_ERROR"
    ErrCodeNotFound   = "NOT_FOUND"
    ErrCodeInternal   = "INTERNAL_ERROR"
)
```

### Context Usage
Use context for structured debugging information:

```go
ctx := map[string]interface{}{
    "user_id": userID,
    "operation": "user_create",
    "timestamp": time.Now(),
}
err := errors.NewWithContext("DATABASE_ERROR", "Insert failed", ctx)
```

### Stack Trace Management
Stack traces are automatically captured when wrapping errors:

```go
// Stack trace captured automatically
wrappedErr := errors.Wrap(originalErr, "SERVICE_ERROR", "Operation failed")

// Access stack trace
if wrappedErr.Stack != nil {
    fmt.Println(wrappedErr.Stack.String())
}
```

## Performance Considerations

- Stack trace capture adds minimal overhead (~1-2μs)
- Error creation is optimized for common cases
- JSON marshaling is efficient with custom implementation
- Memory usage is minimal for typical error scenarios

## Thread Safety

All public APIs are thread-safe for concurrent access. Error instances should not be modified after creation to ensure thread safety. 


---

go-errors • an AGILira library