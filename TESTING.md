# Testing Guide for LazyDnD

This document describes the test suite and how to run tests for LazyDnD.

## Test Coverage

### Current Coverage
- **Overall**: 15.0% of statements
- **panels package**: 22.6% (dice roller and initiative tracker logic)
- **ui package**: 9.7% (handlers, helpers, and state management)

### Test Structure
All tests are located in the `tests/` directory, separate from the source code:
- `tests/dice_roller_test.go` - Tests for dice rolling and initiative functionality
- `tests/save_manager_test.go` - Tests for save/load functionality

This structure keeps tests isolated from production code while maintaining access to all package functionality.

## Running Tests

### Using Makefile (Recommended)

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage report
make test-coverage

# Generate HTML coverage report
make test-coverage-html

# Run all quality checks (format, lint, test)
make check
```

### Using Go Commands Directly

```bash
# Run all tests
go test ./tests/...

# Run tests with verbose output
go test -v ./tests/...

# Run tests with coverage (including all packages)
go test ./tests/... -cover -coverpkg=./...

# Generate coverage report
go test ./tests/... -coverprofile=coverage.out -coverpkg=./...
go tool cover -html=coverage.out
```

## Test Structure

### Dice Roller Tests (`tests/dice_roller_test.go`)

**TestRollDice** - Tests basic dice rolling:
- Simple dice rolls (d20, d6)
- Multiple dice (2d6)
- Dice with modifiers (1d20+5, 1d20-2)
- Invalid commands
- Advantage/disadvantage rolls

**TestRollDiceResults** - Tests dice result ranges:
- Verifies results are within valid ranges
- Tests randomness across multiple rolls

**TestRollInitiative** - Tests initiative rolling:
- Verifies results are between 1-20
- Tests across multiple rolls

**TestParseInput** - Tests input parsing:
- Valid and invalid player names
- Initiative values (positive/negative)
- HP values and changes
- Monster initiative rolls (including "r" for roll)

### Save Manager Tests (`tests/save_manager_test.go`)

**TestSaveState** - Tests state serialization:
- JSON marshaling/unmarshaling
- All fields preserved correctly

**TestSaveLoadCycle** - Tests full save/load:
- Writing to file
- Reading from file
- Data integrity checks

**TestListSaves** - Tests listing save files:
- Finding .json files
- Filtering non-json files

**TestEmptySave** - Tests saving empty state:
- Empty initiative list
- Unstarted combat (turn -1, round 0)

**TestSaveWithMonsterData** - Tests monster saves:
- Monster instance numbers preserved
- Monster metadata preserved

## Writing New Tests

### Test File Naming
- Test files should be named `*_test.go`
- Place test files in the same package as the code being tested

### Test Function Naming
- Test functions must start with `Test`
- Use descriptive names: `TestRollDice`, `TestParseInput`

### Table-Driven Tests
We use table-driven tests for better organization:

```go
func TestExampleFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "Valid input",
            input:   "test",
            want:    "expected output",
            wantErr: false,
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := ExampleFunction(tt.input)
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Best Practices
1. **Test one thing at a time** - Each test should verify a single behavior
2. **Use descriptive names** - Test names should clearly indicate what's being tested
3. **Test edge cases** - Include boundary values, empty inputs, invalid inputs
4. **Keep tests independent** - Tests should not depend on each other
5. **Use subtests** - Use `t.Run()` for table-driven tests

## Test Coverage Goals

### Short-term (Current Sprint)
- ✅ Basic dice rolling functionality
- ✅ Initiative tracker helper functions
- ✅ Navigation handlers (turn tracking, reset combat, scrolling)
- ✅ Save/load functionality
- ✅ Input handlers (Enter, Escape, Backspace, dice history)

### Medium-term
- Action popup logic
- Spell search functionality
- Monster search functionality
- Layout and rendering tests
- More edge cases for existing handlers

### Long-term
- Full UI handler coverage (>80%)
- Integration tests
- End-to-end tests
- Performance benchmarks

## Continuous Integration

Tests are automatically run on:
- Every pull request
- Every push to main branch
- Before releases

## Troubleshooting

### Tests Failing Locally
1. Ensure you have Go 1.19 or higher: `go version`
2. Update dependencies: `go mod download`
3. Clean and rebuild: `make clean && make build`
4. Run tests with verbose output: `make test-verbose`

### Coverage Not Generating
1. Ensure you have write permissions in the directory
2. Check that `go tool cover` is available: `go tool cover -h`
3. Try generating manually: `go test ./... -coverprofile=coverage.out`

## Contributing Tests

When contributing new features:
1. Write tests for new functionality
2. Ensure existing tests pass: `make test`
3. Run coverage check: `make test-coverage`
4. Aim for at least 70% coverage of new code
5. Include tests in your pull request

For more information on contributing, see [CONTRIBUTING.md](./CONTRIBUTING.md).
