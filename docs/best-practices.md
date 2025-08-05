# Best Practices

This document outlines recommended patterns and practices for using go-errors effectively in production applications.

## Error Code Management

### Define Error Codes as Constants
Create a centralized location for error code definitions to ensure consistency across your application.

```go
// errors/codes.go
package errors

const (
    // Validation errors
    ErrCodeValidation     = "VALIDATION_ERROR"
    ErrCodeRequiredField  = "REQUIRED_FIELD"
    ErrCodeInvalidFormat  = "INVALID_FORMAT"
    
    // Database errors
    ErrCodeDatabase       = "DATABASE_ERROR"
    ErrCodeNotFound       = "NOT_FOUND"
    ErrCodeDuplicate      = "DUPLICATE_ENTRY"
    ErrCodeConstraint     = "CONSTRAINT_VIOLATION"
    
    // Network errors
    ErrCodeNetwork        = "NETWORK_ERROR"
    ErrCodeTimeout        = "TIMEOUT"
    ErrCodeConnection     = "CONNECTION_FAILED"
    
    // Authentication errors
    ErrCodeAuth           = "AUTHENTICATION_ERROR"
    ErrCodeUnauthorized   = "UNAUTHORIZED"
    ErrCodeForbidden      = "FORBIDDEN"
    
    // Business logic errors
    ErrCodeBusiness       = "BUSINESS_ERROR"
    ErrCodeInsufficient   = "INSUFFICIENT_RESOURCES"
    ErrCodeQuotaExceeded  = "QUOTA_EXCEEDED"
)
```

### Use Hierarchical Error Codes
Structure error codes to support hierarchical categorization for better error handling.

```go
const (
    // Base categories
    ErrCodeValidation = "VALIDATION_ERROR"
    ErrCodeDatabase   = "DATABASE_ERROR"
    
    // Specific subcategories
    ErrCodeValidationEmail    = "VALIDATION_ERROR.EMAIL"
    ErrCodeValidationPassword = "VALIDATION_ERROR.PASSWORD"
    ErrCodeDatabaseConnection = "DATABASE_ERROR.CONNECTION"
    ErrCodeDatabaseQuery      = "DATABASE_ERROR.QUERY"
)
```

## Error Creation Patterns

### Use Appropriate Constructors
Choose the constructor that best fits your use case:

```go
// Basic error
err := errors.New(ErrCodeValidation, "Invalid input")

// Field validation error
err := errors.NewWithField(ErrCodeValidation, "Email format invalid", "email", email)

// Error with context
ctx := map[string]interface{}{
    "user_id": userID,
    "operation": "user_create",
}
err := errors.NewWithContext(ErrCodeDatabase, "Insert failed", ctx)
```

### Provide Meaningful Messages
Write technical messages that help developers understand and debug issues:

```go
// Good: Specific and actionable
err := errors.New(ErrCodeDatabase, "Failed to insert user: connection timeout after 5s")

// Bad: Too generic
err := errors.New(ErrCodeDatabase, "Database error")
```

### Add User-Friendly Messages
Always provide user-friendly messages for API responses:

```go
err := errors.New(ErrCodeValidation, "Email format invalid").
    WithUserMessage("Please enter a valid email address")
```

## Error Wrapping Strategy

### Wrap at Service Boundaries
Wrap errors when crossing service boundaries to add context:

```go
func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    err := s.repo.Create(ctx, user)
    if err != nil {
        return errors.Wrap(err, ErrCodeDatabase, "Failed to create user in database")
    }
    return nil
}
```

### Preserve Original Errors
Always preserve the original error in the cause chain:

```go
// Good: Preserves original error
dbErr := sql.ErrNoRows
wrappedErr := errors.Wrap(dbErr, ErrCodeNotFound, "User not found")

// Bad: Loses original error
err := errors.New(ErrCodeNotFound, "User not found")
```

### Add Context Progressively
Add context as errors bubble up through the call stack:

```go
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    user, err := h.service.CreateUser(r.Context(), userData)
    if err != nil {
        // Add HTTP-specific context
        httpErr := errors.Wrap(err, ErrCodeHTTP, "Failed to process create user request")
        http.Error(w, httpErr.UserMessage(), http.StatusBadRequest)
        return
    }
}
```

## Error Handling Patterns

### Use Error Codes for Programmatic Handling
Leverage error codes for conditional logic:

```go
func handleError(err error) {
    switch {
    case errors.HasCode(err, ErrCodeValidation):
        // Handle validation errors
        logValidationError(err)
        return validationResponse(err)
        
    case errors.HasCode(err, ErrCodeDatabase):
        // Handle database errors
        logDatabaseError(err)
        return internalServerError()
        
    case errors.HasCode(err, ErrCodeNotFound):
        // Handle not found errors
        return notFoundResponse(err)
        
    default:
        // Handle unknown errors
        logUnknownError(err)
        return internalServerError()
    }
}
```

### Implement Retry Logic
Use the Retryable interface for implementing retry mechanisms:

```go
func (c *Client) makeRequest(req *Request) (*Response, error) {
    for attempts := 0; attempts < maxRetries; attempts++ {
        resp, err := c.doRequest(req)
        if err == nil {
            return resp, nil
        }
        
        // Check if error is retryable
        if retryable, ok := err.(errors.Retryable); ok && retryable.IsRetryable() {
            time.Sleep(backoff(attempts))
            continue
        }
        
        return nil, err
    }
    return nil, errors.New(ErrCodeTimeout, "Max retries exceeded")
}
```

### Log Errors Appropriately
Log errors with appropriate detail levels:

```go
func logError(err error) {
    if apiErr, ok := err.(*errors.Error); ok {
        // Log structured error information
        log.Printf("Error: %s, Code: %s, Context: %+v", 
            apiErr.Message, apiErr.Code, apiErr.Context)
        
        // Log stack trace for debugging
        if apiErr.Stack != nil {
            log.Printf("Stack trace: %s", apiErr.Stack.String())
        }
    } else {
        // Log standard errors
        log.Printf("Standard error: %v", err)
    }
}
```

## API Integration

### Return Structured Errors in APIs
Use structured errors for consistent API responses:

```go
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    user, err := h.service.CreateUser(r.Context(), userData)
    if err != nil {
        var apiErr *errors.Error
        if errors.As(err, &apiErr) {
            // Return structured error response
            response := ErrorResponse{
                Code:    apiErr.Code,
                Message: apiErr.UserMessage(),
                Details: apiErr.Context,
            }
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(getHTTPStatus(apiErr.Code))
            json.NewEncoder(w).Encode(response)
        } else {
            // Fallback for non-structured errors
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
```

### Map Error Codes to HTTP Status
Create a mapping function for error codes to HTTP status codes:

```go
func getHTTPStatus(code errors.ErrorCode) int {
    switch code {
    case ErrCodeValidation, ErrCodeRequiredField, ErrCodeInvalidFormat:
        return http.StatusBadRequest
    case ErrCodeNotFound:
        return http.StatusNotFound
    case ErrCodeUnauthorized:
        return http.StatusUnauthorized
    case ErrCodeForbidden:
        return http.StatusForbidden
    case ErrCodeDuplicate:
        return http.StatusConflict
    case ErrCodeTimeout, ErrCodeNetwork:
        return http.StatusGatewayTimeout
    default:
        return http.StatusInternalServerError
    }
}
```

## Testing Patterns

### Test Error Codes
Verify that errors have the correct codes:

```go
func TestCreateUser_ValidationError(t *testing.T) {
    _, err := service.CreateUser(ctx, invalidUser)
    
    if err == nil {
        t.Fatal("Expected error, got nil")
    }
    
    if !errors.HasCode(err, ErrCodeValidation) {
        t.Errorf("Expected validation error, got: %v", err)
    }
    
    // Check user message
    if apiErr, ok := err.(*errors.Error); ok {
        if apiErr.UserMessage() == "" {
            t.Error("Expected user message to be set")
        }
    }
}
```

### Test Error Wrapping
Ensure errors are properly wrapped:

```go
func TestService_CreateUser_WrapsDatabaseError(t *testing.T) {
    // Mock database to return error
    mockDB.EXPECT().Create(gomock.Any(), gomock.Any()).Return(sql.ErrNoRows)
    
    _, err := service.CreateUser(ctx, user)
    
    if err == nil {
        t.Fatal("Expected error, got nil")
    }
    
    // Check that original error is preserved
    root := errors.RootCause(err)
    if root != sql.ErrNoRows {
        t.Errorf("Expected root cause to be sql.ErrNoRows, got: %v", root)
    }
    
    // Check that error is wrapped with service context
    if !errors.HasCode(err, ErrCodeDatabase) {
        t.Error("Expected database error code")
    }
}
```

## Performance Considerations

### Minimize Stack Trace Overhead
Stack traces are automatically captured when wrapping errors. Consider the performance impact:

```go
// Only wrap errors that need debugging context
if debugMode {
    err = errors.Wrap(err, ErrCodeInternal, "Additional context")
}
```

### Reuse Error Instances
For frequently occurring errors, consider creating reusable error instances:

```go
var (
    ErrUserNotFound = errors.New(ErrCodeNotFound, "User not found").
        WithUserMessage("User not found")
    ErrInvalidEmail = errors.New(ErrCodeValidation, "Invalid email format").
        WithUserMessage("Please enter a valid email address")
)
```

### Optimize Context Usage
Use context sparingly to avoid memory overhead:

```go
// Good: Minimal context
ctx := map[string]interface{}{
    "user_id": userID,
}

// Bad: Excessive context
ctx := map[string]interface{}{
    "user_id": userID,
    "timestamp": time.Now(),
    "request_id": requestID,
    "session_id": sessionID,
    "ip_address": ipAddress,
    // ... many more fields
}
```

## Security Considerations

### Sanitize Error Messages
Never expose sensitive information in error messages:

```go
// Good: Sanitized message
err := errors.New(ErrCodeDatabase, "Database operation failed")

// Bad: Exposes sensitive information
err := errors.New(ErrCodeDatabase, "Failed to connect to database: password=secret123")
```

### Validate Error Codes
Ensure error codes are valid and consistent:

```go
func validateErrorCode(code errors.ErrorCode) error {
    validCodes := map[errors.ErrorCode]bool{
        ErrCodeValidation: true,
        ErrCodeDatabase:   true,
        ErrCodeNetwork:    true,
        // ... other valid codes
    }
    
    if !validCodes[code] {
        return errors.New(ErrCodeInternal, "Invalid error code")
    }
    return nil
}
```

## Monitoring and Observability

### Track Error Metrics
Monitor error rates and patterns:

```go
func trackError(err error) {
    if apiErr, ok := err.(*errors.Error); ok {
        // Increment error counter by code
        errorCounter.WithLabelValues(string(apiErr.Code)).Inc()
        
        // Track error context
        if apiErr.Context != nil {
            for key, value := range apiErr.Context {
                errorContextGauge.WithLabelValues(key).Set(1)
            }
        }
    }
}
```

### Structured Logging
Use structured logging for better error analysis:

```go
func logStructuredError(err error) {
    if apiErr, ok := err.(*errors.Error); ok {
        log.WithFields(log.Fields{
            "error_code":    apiErr.Code,
            "error_message": apiErr.Message,
            "user_message":  apiErr.UserMessage,
            "timestamp":     apiErr.Timestamp,
            "context":       apiErr.Context,
            "retryable":     apiErr.Retryable,
        }).Error("Application error occurred")
    }
}
``` 

---

go-errors • an AGILira library