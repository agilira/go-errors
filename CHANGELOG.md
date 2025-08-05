# Changelog

All notable changes to the go-errors library will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive test suite achieving 100% code coverage
- Edge case testing for all error handling functions
- Thread safety validation with race condition detection
- Security analysis with gosec
- Code quality checks with go vet, golint, and go fmt

### Changed
- Updated README.md with accurate 100% coverage badge
- Enhanced code reliability through extensive testing

### Fixed
- Edge cases in JSON marshaling with nil stack traces
- User message handling edge cases
- Stack trace formatting edge cases
- **Critical**: Fixed stack trace line number formatting bug (was converting line numbers to Unicode characters)
- Simplified RootCause function to use standard library errors.Unwrap
- Removed duplicate interface method implementations from test file
- Added proper Godoc comments for As method behavior
- Implemented missing ErrorCode() and IsRetryable() methods for interface compliance

### Security
- No security vulnerabilities detected by gosec
- Thread-safe error operations confirmed

## [1.0.0] - Initial Release

### Added
- Core error handling library with structured error types
- Error codes support for categorization
- Context and field validation support
- Stack trace capture and formatting
- JSON serialization for API/microservices
- User-friendly message support
- Retryable error interface
- Helper functions for error wrapping and chain traversal
- Standard library compatibility (errors.Is, errors.As, errors.Unwrap)
- Comprehensive documentation

### Features
- `New()` - Create basic errors with code and message
- `NewWithField()` - Create errors with field validation context
- `NewWithContext()` - Create errors with additional context
- `Wrap()` - Wrap existing errors with new context and stack trace
- `RootCause()` - Find the original error in a chain
- `HasCode()` - Search for specific error codes in a chain
- `Is()` - Standard library compatibility for error comparison
- `As()` - Standard library compatibility for error type assertion
- `Error()` - Implement error interface with formatted message
- `Unwrap()` - Return underlying cause for error chaining
- `MarshalJSON()` - Custom JSON serialization
- `WithUserMessage()` - Set user-friendly error messages
- `UserMessage()` - Get appropriate message for users
- `CaptureStacktrace()` - Capture current call stack
- `String()` - Format stack trace as human-readable string

### Technical Details
- Zero external dependencies
- Go 1.24.4+ compatibility
- MPL-2.0 license
- Cross-platform support (Windows, Linux, macOS)
- High performance with minimal memory overhead
- Thread-safe operations 