# go-errors: Structured, contextual error handling for Go
### an AGILira library

go-errors is a fast, structured, and context-aware error handling library for Go.
Originally built for [Orpheus](https://github.com/agilira/orpheus), it provides error codes, stack traces, user messages, and JSON support with near-zero overhead through [Timecache](https://github.com/agilira/go-timecache) integration.

[![CI](https://github.com/agilira/go-errors/actions/workflows/ci.yml/badge.svg)](https://github.com/agilira/go-errors/actions/workflows/ci.yml)
[![Security](https://img.shields.io/badge/Security-gosec-brightgreen)](https://github.com/agilira/go-errors/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/agilira/go-errors?v=2)](https://goreportcard.com/report/github.com/agilira/go-errors)
[![Coverage](https://codecov.io/gh/agilira/go-errors/branch/main/graph/badge.svg)](https://codecov.io/gh/agilira/go-errors)
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/agilira/go-errors.svg)](https://pkg.go.dev/github.com/agilira/go-errors)

**[Features](#features) • [Performance](#performance) • [Quick Start](#quick-start) • [JSON Output](#json-output) • [Testing & Coverage](#testing--coverage) • [Documentation](#documentation)**

## Features
- Structured error type: code, message, context, cause, severity
- Stacktrace support (optional, lightweight)
- User and technical messages
- Custom error codes (user-defined)
- JSON serialization for API/microservices
- Retryable and interface-based error handling
- Helpers for wrapping, root cause, code search
- Modular, fully tested, high coverage

## Compatibility and Support

go-plugins is designed for Go 1.23+ environments and follows Long-Term Support guidelines to ensure consistent performance across production deployments.

## Performance

```
AMD Ryzen 5 7520U with Radeon Graphics
BenchmarkNew-8                          12971414      86.61 ns/op    208 B/op    2 allocs/op
BenchmarkNewWithField-8                 12536019      87.98 ns/op    208 B/op    2 allocs/op
BenchmarkNewWithContext-8               20815206      57.17 ns/op    160 B/op    1 allocs/op
BenchmarkWrap-8                          2111182     558.0 ns/op     264 B/op    4 allocs/op
BenchmarkMethodChaining-8                5201632     220.8 ns/op     504 B/op    3 allocs/op
BenchmarkHasCode-8                      325451757       3.66 ns/op      0 B/op    0 allocs/op
BenchmarkRootCause-8                    144666518       8.18 ns/op      0 B/op    0 allocs/op
BenchmarkMarshalJSON-8                    449632     2603 ns/op      568 B/op    7 allocs/op
```

**Reproduce benchmarks**:
```bash
go test -bench=. -benchmem
```

## Quick Start

### Installation

```sh
go get github.com/agilira/go-errors
```

### Quick Example
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

## JSON Output
Errors automatically serialize to structured JSON:
```go
err := errors.New("VALIDATION_ERROR", "Email format invalid").
    WithUserMessage("Please enter a valid email address").
    WithContext("field", "email").
    WithSeverity("warning")

jsonData, _ := json.Marshal(err)
// Output: {"code":"VALIDATION_ERROR","message":"Email format invalid","user_msg":"Please enter a valid email address","context":{"field":"email"},"severity":"warning","timestamp":"2025-01-27T10:30:00Z",...}
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
