# Tests Directory

This directory contains all unit tests for the LazyDnD project, isolated from the source code for better organization and maintainability.

## Structure

All test files are in the `tests` package and import the packages they test:

- `dice_roller_test.go` - Tests for dice rolling and initiative logic
- `save_manager_test.go` - Tests for save/load functionality

## Running Tests

### Quick Start
```bash
# From project root
make test

# With verbose output
make test-verbose

# With coverage
make test-coverage
```

### Direct Go Commands
```bash
# Run all tests
go test ./tests/...

# Run with coverage for all packages
go test ./tests/... -cover -coverpkg=./...

# Generate HTML coverage report
go test ./tests/... -coverprofile=coverage.out -coverpkg=./...
go tool cover -html=coverage.out
```

## Adding New Tests

1. Create a new test file in this directory: `tests/your_feature_test.go`
2. Use package `tests`
3. Import the packages you need to test:
   ```go
   package tests

   import (
       "testing"
       "lazydnd/ui"
       "lazydnd/panels"
   )
   ```
4. Write your test functions following Go testing conventions
5. Run `make test` to verify

## Test Organization

Tests are organized by feature rather than by source package:
- **Functionality-based**: Tests focus on specific features (dice rolling, saving, etc.)
- **Integration-friendly**: Tests can easily test interactions between packages
- **Clean separation**: Production code stays separate from test code

## Coverage

Current test coverage: **8.6% of statements**

To improve coverage:
1. Identify untested functions: `make test-coverage`
2. Add tests for critical paths first
3. Focus on exported functions and public APIs
4. Test edge cases and error conditions

## Best Practices

- Use table-driven tests for multiple test cases
- Test one behavior per test function
- Use descriptive test names: `TestFeature_Scenario_ExpectedBehavior`
- Keep tests independent - no shared state
- Mock external dependencies when needed
- Test both success and failure cases

For more details, see the main [TESTING.md](../TESTING.md) in the project root.
