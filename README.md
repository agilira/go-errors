# go-errors: Structured, contextual error handling for Go
### an AGILira library

go-errors is a fast, structured, and context-aware error handling library for Go.
Built for Styx, it includes error codes, stack traces, user messages, and JSON support — with zero overhead.

[![CI](https://github.com/agilira/go-errors/actions/workflows/ci.yml/badge.svg)](https://github.com/agilira/go-errors/actions/workflows/ci.yml)
[![Security](https://img.shields.io/badge/Security-gosec-brightgreen)](https://github.com/agilira/go-errors/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)](https://github.com/agilira/go-errors/actions/workflows/ci.yml)

## Features
- Structured error type: code, message, context, cause, severity
- Stacktrace support (optional, lightweight)
- User and technical messages
- Custom error codes (user-defined)
- JSON serialization for API/microservices
- Retryable and interface-based error handling
- Helpers for wrapping, root cause, code search
- 100% Go standard library, no external dependencies
- Modular, fully tested, high coverage

## Installation
```sh
go get github.com/agilira/go-errors
```

## Quick Example
```go
import "github.com/agilira/go-errors"

const ErrCodeValidation = "VALIDATION_ERROR"

func validateUser(username string) error {
    if username == "" {
        return errors.New(ErrCodeValidation, "Username is required").WithUserMessage("Please enter a username.")
    }
    return nil
}
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

## Documentation

Comprehensive documentation is available in the [docs](./docs/) folder:

- **[API Reference](./docs/api.md)** - Complete API documentation with examples
- **[Usage Guide](./docs/usage.md)** - Getting started and basic usage patterns
- **[Best Practices](./docs/best-practices.md)** - Production-ready patterns and recommendations
- **[Integration Guide](./docs/integration.md)** - Migration guide and advanced integration patterns

---

go-errors • an AGILira library
