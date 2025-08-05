# Usage Guide

This guide covers all core features of go-errors, with practical examples.

## Getting Started

### Installation
```sh
go get github.com/AGILira/go-errors
```

### Import
```go
import goerrors "github.com/AGILira/go-errors"
```

### Creating and Wrapping Errors
```go
const ErrCodeValidation = "VALIDATION_ERROR"
err := goerrors.New(ErrCodeValidation, "Validation failed")
wrapped := goerrors.Wrap(err, ErrCodeValidation, "Additional context")
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
err := goerrors.New("VALIDATION_ERROR", "Invalid input").WithUserMessage("Please check your input.")
msg := err.UserMessage() // User message if set, else technical message
```

## Stacktrace
```go
err := goerrors.Wrap(errors.New("fail"), "CODE", "Context")
if err.Stack != nil {
    fmt.Println(err.Stack.String())
}
```

## Helpers
```go
root := goerrors.RootCause(err)
if goerrors.HasCode(err, "VALIDATION_ERROR") { /* ... */ }
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
err := goerrors.New("DB_ERROR", "Query failed").WithUserMessage("Database unavailable").WithContext("query", "SELECT * FROM users")
err = goerrors.Wrap(err, "SERVICE_ERROR", "Service layer failure")
```

### Retryable Errors
```go
err := goerrors.New("TEMP_ERROR", "Temporary failure")
err.Retryable = true
if retry, ok := interface{}(err).(goerrors.Retryable); ok && retry.IsRetryable() {
    // Retry logic
}
```

### API Integration
```go
// In a REST handler
if err != nil {
    apiErr, _ := err.(*goerrors.Error)
    http.Error(w, apiErr.UserMessage(), 400)
}
``` 

---

go-errors • an AGILira library