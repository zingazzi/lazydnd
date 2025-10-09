# Contributing to LazyD&D

Thank you for your interest in contributing to LazyD&D! We welcome contributions of all kinds, including bug fixes, new features, documentation improvements, and suggestions.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [How to Contribute](#how-to-contribute)
- [Coding Guidelines](#coding-guidelines)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)
- [Reporting Bugs](#reporting-bugs)
- [Feature Requests](#feature-requests)

## Code of Conduct

This project follows a simple code of conduct: be respectful, constructive, and helpful. We're all here to make a great tool for D&D enthusiasts.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/lazydnd.git
   cd lazydnd
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/zingazzi/lazydnd.git
   ```

## Development Setup

### Prerequisites

- Go 1.19 or higher
- Git
- A terminal with color support

### Building from Source

```bash
# Install dependencies
go mod download

# Build the project
go build -o lazydnd

# Run the application
./lazydnd
```

### Running Tests

All tests are located in the `tests/` directory. Run them using:

```bash
# Run all tests
make test

# Run with verbose output
make test-verbose

# Run with coverage report
make test-coverage

# Or use Go directly
go test ./tests/...
```

**Important**: All tests must pass before submitting a pull request.

## Project Structure

Please review [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md) for detailed information about the codebase organization. Key directories:

- `main.go` - Application entry point
- `ui/` - User interface components and handlers
  - `model.go` - Main Bubble Tea model
  - `handlers_*.go` - Key event handlers organized by category
  - `layout_*.go` - View rendering functions
  - `*_popup.go` - Popup dialogs
- `panels/` - Panel-specific logic (dice, initiative, spells, monsters)
- `assets/` - JSON data files for spells and monsters

## How to Contribute

### 1. Create a Branch

Create a feature branch from `main`:

```bash
git checkout -b feature/your-feature-name
```

Use descriptive branch names:
- `feature/add-spell-filtering`
- `bugfix/initiative-sort-issue`
- `docs/update-installation`

### 2. Make Your Changes

- Write clean, readable code
- Follow the existing code style
- Add comments where the operation isn't clear
- Include file path/name as a one-line comment at the top of new files
- **Add or update tests** for your changes (see [Testing Guidelines](#testing-guidelines))
- Test your changes thoroughly

### 3. Commit Your Changes

```bash
git add .
git commit -m "Add feature: description of your changes"
```

See [Commit Messages](#commit-messages) for guidelines.

### 4. Push to Your Fork

```bash
git push origin feature/your-feature-name
```

### 5. Open a Pull Request

Go to the [LazyD&D repository](https://github.com/zingazzi/lazydnd) and open a Pull Request from your fork.

## Coding Guidelines

### General Principles

- **Modularity**: Keep functions focused and single-purpose
- **DRY**: Don't repeat yourself - extract common logic
- **Performance**: Consider efficiency, especially for large datasets
- **Security**: Validate inputs and handle errors gracefully

### Go-Specific Guidelines

- Use **ES module syntax** where applicable
- Follow **Go conventions**: use `gofmt` for formatting
- **Error handling**: Always handle errors explicitly
- **Variable naming**: Use descriptive names (e.g., `monsterIndex` not `mi`)
- **Comments**: Describe purpose, not effect
  - Good: `// Calculate panel dimensions based on terminal size`
  - Bad: `// Set width to termWidth / 2`

### Code Style

- Use tabs for indentation (Go standard)
- Keep functions under 50 lines when possible
- Extract complex logic into helper functions
- Use struct types for related data (see `PanelDimensions`, `ContentHeights`)
- Organize handler functions using maps instead of large switch statements

### Example: Handler Pattern

```go
// Good: Map-based handler pattern
var keyHandlers = map[string]func(*Model) tea.Cmd{
    "q":     handleQuit,
    "tab":   handleTabForward,
    "enter": handleEnter,
}

// Bad: Large switch statement
switch key {
case "q":
    return handleQuit(m)
case "tab":
    return handleTabForward(m)
// ... 50 more cases
}
```

## Commit Messages

Write clear, concise commit messages:

### Format

```
<type>: <subject>

<body (optional)>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, no logic change)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples

```
feat: add spell level filtering to spells panel

fix: initiative tracker sort order after HP edit

docs: update installation instructions for ARM64

refactor: extract scroll logic into separate functions
```

## Pull Request Process

1. **Update documentation** if you've changed functionality
2. **Add or update tests** for your changes
3. **Ensure all tests pass** - run `make test` before submitting
4. **Test thoroughly** - ensure the app runs without crashes
5. **Keep PRs focused** - one feature or fix per PR
6. **Describe your changes** clearly in the PR description
7. **Link related issues** using keywords (e.g., "Fixes #123")
8. **Be responsive** to feedback and requested changes

### PR Description Template

```markdown
## Description
Brief description of what this PR does.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Refactoring

## Testing
- [ ] All existing tests pass (`make test`)
- [ ] Added new tests for new functionality
- [ ] Manually tested the changes
- [ ] Tested edge cases

How did you test these changes?

## Related Issues
Fixes #(issue number)

## Screenshots (if applicable)
```

## Reporting Bugs

Found a bug? [Open an issue](https://github.com/zingazzi/lazydnd/issues/new) with:

1. **Clear title** describing the issue
2. **Steps to reproduce** the bug
3. **Expected behavior** vs **actual behavior**
4. **Environment details**:
   - OS and version
   - Terminal emulator
   - LazyD&D version
5. **Screenshots or logs** if applicable

### Bug Report Template

```markdown
**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
1. Go to '...'
2. Press '...'
3. See error

**Expected behavior**
What you expected to happen.

**Environment:**
- OS: [e.g., macOS 13.0, Ubuntu 22.04]
- Terminal: [e.g., iTerm2, GNOME Terminal]
- LazyD&D version: [e.g., v1.0.0]

**Screenshots**
If applicable, add screenshots.
```

## Feature Requests

Have an idea? [Open an issue](https://github.com/zingazzi/lazydnd/issues/new) with:

1. **Clear title** describing the feature
2. **Use case** - why is this feature needed?
3. **Proposed solution** - how should it work?
4. **Alternatives considered** (optional)

### Feature Request Template

```markdown
**Is your feature request related to a problem?**
A clear description of the problem.

**Describe the solution you'd like**
What you want to happen.

**Describe alternatives you've considered**
Other solutions you've thought about.

**Additional context**
Any other context or screenshots.
```

## Testing Guidelines

### Writing Tests

All tests should be placed in the `tests/` directory and use the `tests` package:

```go
// tests/your_feature_test.go
package tests

import (
    "testing"
    "lazydnd/ui"
    "lazydnd/panels"
)

func TestYourFeature(t *testing.T) {
    // Your test code here
}
```

### Test Requirements

When contributing code, please ensure:

1. **All existing tests pass** - Run `make test` before submitting
2. **Add tests for new features** - New functionality should include tests
3. **Add tests for bug fixes** - Add a test that would have caught the bug
4. **Test edge cases** - Consider boundary conditions, empty inputs, invalid data
5. **Use table-driven tests** - For multiple test cases of the same function

### Example: Table-Driven Test

```go
func TestRollDice(t *testing.T) {
    tests := []struct {
        name         string
        command      string
        wantContains string
        wantError    bool
    }{
        {
            name:         "Simple d20 roll",
            command:      "1d20",
            wantContains: "d20:",
            wantError:    false,
        },
        {
            name:         "Invalid command",
            command:      "invalid",
            wantContains: "Invalid",
            wantError:    true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := panels.RollDice(tt.command)
            if !strings.Contains(result, tt.wantContains) {
                t.Errorf("got %q, want to contain %q", result, tt.wantContains)
            }
        })
    }
}
```

### Running Specific Tests

```bash
# Run all tests
make test

# Run with verbose output
make test-verbose

# Run specific test
go test ./tests/... -run TestRollDice

# Run with coverage
make test-coverage
```

### Test Coverage

- Current coverage: **8.6%** of statements
- Goal: Aim for **70%+** coverage for new code
- Focus on testing:
  - Exported functions and public APIs
  - Critical game logic (dice rolling, initiative, save/load)
  - Edge cases and error conditions

See [TESTING.md](./TESTING.md) for more detailed testing documentation.

## Development Tips

### Testing UI Changes

- Test in different terminal sizes
- Verify scrolling behavior with long content
- Check that keybindings work as expected
- Test edge cases (empty lists, very long names, etc.)

### Working with Bubble Tea

- Understand the Elm Architecture: Model → Update → View
- Commands (`tea.Cmd`) are for side effects
- Messages (`tea.Msg`) trigger updates
- Keep the `Update` function pure when possible

### Adding New Panels

1. Create panel logic in `panels/`
2. Add panel state to `Model` in `ui/model.go`
3. Add rendering logic in `ui/layout_content.go`
4. Add keybindings in appropriate `ui/handlers_*.go` file
5. Update help text in `ui/help_popup.go`

### Debugging

- Use `fmt.Fprintf(os.Stderr, "debug: %v\n", value)` for debugging (stdout is used by Bubble Tea)
- Test with `go run main.go`
- Check for panics and handle errors gracefully

## Questions?

If you have questions about contributing:

- Open a [discussion](https://github.com/zingazzi/lazydnd/discussions)
- Comment on relevant issues
- Reach out via the project's communication channels

## Thank You!

Your contributions help make LazyD&D better for the entire D&D community. Whether it's code, documentation, bug reports, or feature ideas - every contribution matters!

Happy coding, and may your rolls be high! 🎲
