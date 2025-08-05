# Integration, Migration & FAQ

## Interfaces

go-errors provides interfaces for advanced error handling:
- **ErrorCoder**: Extracts the error code
- **Retryable**: Indicates if an error is retryable
- **UserMessager**: Extracts a user-friendly message

### Example
```go
var coder goerrors.ErrorCoder = err
code := coder.ErrorCode()

var retry goerrors.Retryable = err
if retry.IsRetryable() { /* ... */ }

var um goerrors.UserMessager = err
msg := um.UserMessage()
```

## Testing & Coverage
Run all tests:
```sh
go test -v ./...
```
Check coverage:
```sh
go test -cover ./...
```
- Write tests for all custom error codes and logic in your application.
- Use table-driven tests for error scenarios.
- Aim for high coverage to ensure reliability.

## Integration in Real Projects
- Use go-errors for all error creation and propagation in your services.
- Wrap errors at API/service boundaries to add context and user messages.
- Serialize errors to JSON for API responses.
- Use error codes for programmatic error handling in clients.

### Example: REST API Handler
```go
func handler(w http.ResponseWriter, r *http.Request) {
    err := doSomething()
    if err != nil {
        apiErr, _ := err.(*goerrors.Error)
        http.Error(w, apiErr.UserMessage(), 400)
        // Optionally log apiErr.Stack.String() for debugging
    }
}
```

## Migration Guide: From Standard Errors to go-errors

This guide helps you migrate your Go project from standard error handling to go-errors, step by step.

### 1. Install go-errors
```sh
go get github.com/AGILira/go-errors
```

### 2. Update Imports
Replace:
```go
import "errors"
```
with:
```go
import goerrors "github.com/AGILira/go-errors"
```

### 3. Define Error Codes
Instead of using only error messages, define string constants for error codes:
```go
const ErrCodeValidation = "VALIDATION_ERROR"
```

### 4. Replace errors.New and fmt.Errorf
**Before:**
```go
return errors.New("validation failed")
```
**After:**
```go
return goerrors.New(ErrCodeValidation, "Validation failed")
```

### 5. Wrapping Errors
**Before:**
```go
return fmt.Errorf("db error: %w", err)
```
**After:**
```go
return goerrors.Wrap(err, "DB_ERROR", "db error")
```

### 6. Add User Messages (Optional)
```go
err := goerrors.New(ErrCodeValidation, "Validation failed").WithUserMessage("Please check your input.")
```

### 7. Use Helpers for Inspection
Replace manual error checks with helpers:
```go
if goerrors.HasCode(err, ErrCodeValidation) { /* ... */ }
root := goerrors.RootCause(err)
```

### 8. Update API/Handler Logic
Return user-friendly messages in APIs:
```go
apiErr, _ := err.(*goerrors.Error)
http.Error(w, apiErr.UserMessage(), 400)
```

### 9. Testing
Update your tests to check error codes and user messages, not just error strings.

#### Example: Before and After
**Before:**
```go
if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        return errors.New("not found")
    }
    return fmt.Errorf("db error: %w", err)
}
```
**After:**
```go
if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        return goerrors.New("NOT_FOUND", "not found").WithUserMessage("Resource not found")
    }
    return goerrors.Wrap(err, "DB_ERROR", "db error")
}
```

### Tips
- Migrate incrementally: start from leaf packages and move up.
- Use go-errors everywhere for consistency.
- Use error codes for programmatic handling, user messages for UI/API.

## FAQ

**Q: Can I define my own error codes?**
A: Yes, define them as string constants in your application.

**Q: Is go-errors compatible with the standard errors package?**
A: Yes, it supports errors.Is, errors.As, and errors.Unwrap.

**Q: How do I attach extra context to an error?**
A: Use the WithContext method.

**Q: How do I return user-friendly messages in APIs?**
A: Use WithUserMessage and UserMessage().

**Q: Is stacktrace capture expensive?**
A: No, it is lightweight and only captured when wrapping errors. 

---

go-errors • an AGILira library