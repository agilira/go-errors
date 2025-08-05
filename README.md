# go-errors
### an AGILira library

Reusable, modular error handling library for Go.

[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://www.mozilla.org/en-US/MPL/2.0/)
[![Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)]()

## Overview
**go-errors** provides robust, structured, and extensible error management for modern Go projects. It supports error codes, context, stacktraces, user messages, JSON serialization, and more—without sacrificing simplicity or performance.

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