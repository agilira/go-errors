# Usage Guide

This guide covers all core features of go-errors, with practical examples and real-world usage patterns.

## Getting Started

### Installation
```sh
go get github.com/agilira/go-errors
```

### Import
```go
import "github.com/agilira/go-errors"
```

### Creating and Wrapping Errors
```go
const ErrCodeValidation = "VALIDATION_ERROR"
err := errors.New(ErrCodeValidation, "Validation failed")
wrapped := errors.Wrap(err, ErrCodeValidation, "Additional context")
```

## Error Structure
- `Code`: Application-defined error code
- `Message`: Technical message
- `UserMsg`: User-friendly message
- `Field`, `Value`: For validation errors
- `Context`: Additional data
- `Timestamp`: Creation time
- `Cause`: Underlying error
- `Severity`: e.g. "error", "warning"
- `Stack`: Stacktrace (if present)
- `Retryable`: Indicates if error is retryable

## User Messages
```go
err := errors.New("VALIDATION_ERROR", "Invalid input").WithUserMessage("Please check your input.")
msg := err.UserMessage() // User message if set, else technical message
```

## Stacktrace
```go
err := errors.Wrap(errors.New("fail"), "CODE", "Context")
if err.Stack != nil {
    fmt.Println(err.Stack.String())
}
```

## Helpers
```go
root := errors.RootCause(err)
if errors.HasCode(err, "VALIDATION_ERROR") { /* ... */ }
if errors.Is(err, targetErr) { /* ... */ }
if errors.As(err, &target) { /* ... */ }
```

## JSON Serialization
```go
b, _ := json.Marshal(err)
fmt.Println(string(b))
```

## Advanced Examples

### Chaining and Context
```go
err := errors.New("DB_ERROR", "Query failed").
    WithUserMessage("Database unavailable").
    WithContext("query", "SELECT * FROM users").
    WithSeverity("critical")
err = errors.Wrap(err, "SERVICE_ERROR", "Service layer failure")
```

### Retryable Errors
```go
err := errors.New("TEMP_ERROR", "Temporary failure").AsRetryable()
if retry, ok := err.(errors.Retryable); ok && retry.IsRetryable() {
    // Retry logic
}
```

### API Integration
```go
// In a REST handler
if err != nil {
    apiErr, _ := err.(*errors.Error)
    http.Error(w, apiErr.UserMessage(), 400)
}
```

### Interface-Based Error Handling
```go
// Type-safe error code extraction
var coder errors.ErrorCoder = err
code := coder.ErrorCode()

// Check if error is retryable
var retry errors.Retryable = err
if retry.IsRetryable() {
    // Implement retry logic
}

// Get user-friendly message
var um errors.UserMessager = err
userMsg := um.UserMessage()
```

### Error Severity Levels
```go
// Different severity levels
err := errors.New("VALIDATION_ERROR", "Invalid input").WithSeverity("warning")
err = errors.New("DATABASE_ERROR", "Connection lost").WithSeverity("critical")
err = errors.New("INFO_ERROR", "Operation completed").WithSeverity("info")
``` 

---

go-errors • an AGILira library