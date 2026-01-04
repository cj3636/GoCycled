# Contributing to GoCycled

Thank you for your interest in contributing to GoCycled! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Building and Testing](#building-and-testing)
- [Coding Standards](#coding-standards)
- [Submitting Changes](#submitting-changes)
- [Reporting Bugs](#reporting-bugs)
- [Feature Requests](#feature-requests)

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/GoCycled.git
   cd GoCycled
   ```
3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/cj3636/GoCycled.git
   ```

## Development Setup

### Prerequisites

- Go 1.18 or later
- Make (optional, for using Makefile)
- Git

### Setting Up Development Environment

1. Install Go dependencies:
   ```bash
   go mod download
   ```

2. Build the project:
   ```bash
   make build
   # or
   go build -o rc ./cmd/rc
   ```

3. Run tests:
   ```bash
   make test
   # or
   go test ./...
   ```

## Building and Testing

### Building

```bash
# Using Make
make build

# Using go directly
go build -o rc ./cmd/rc

# Build and install
make install
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-cover

# Run tests for specific package
go test ./pkg/config
go test ./pkg/trash
```

### Formatting Code

```bash
# Format all Go files
make fmt

# Or use go fmt directly
go fmt ./...
```

### Linting

```bash
# Run linter (requires golangci-lint)
make lint
```

## Coding Standards

### Go Code Style

- Follow standard Go conventions and idioms
- Use `go fmt` to format code
- Write clear, descriptive variable and function names
- Add comments for exported functions and types
- Keep functions small and focused

### Testing

- Write unit tests for new functionality
- Ensure all tests pass before submitting
- Aim for good test coverage
- Use table-driven tests where appropriate

### Commit Messages

- Use clear, descriptive commit messages
- Start with a verb in the imperative mood (e.g., "Add", "Fix", "Update")
- Keep the first line under 72 characters
- Add detailed description if needed

Example:
```
Add restore command with interactive selection

- Implement interactive file selection using basic UI
- Add support for restoring by original path
- Include tests for restore functionality
```

## Submitting Changes

1. Create a new branch for your feature or bugfix:
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. Make your changes and commit:
   ```bash
   git add .
   git commit -m "Add new feature"
   ```

3. Run tests and ensure they pass:
   ```bash
   make test
   ```

4. Push to your fork:
   ```bash
   git push origin feature/my-new-feature
   ```

5. Create a Pull Request on GitHub with:
   - Clear title describing the change
   - Description of what was changed and why
   - Any related issue numbers
   - Test results

### Pull Request Guidelines

- Keep PRs focused on a single feature or fix
- Update documentation if needed
- Add tests for new functionality
- Ensure all tests pass
- Follow the existing code style
- Respond to review feedback promptly

## Reporting Bugs

When reporting bugs, please include:

1. **Description**: Clear description of the bug
2. **Steps to Reproduce**: Detailed steps to reproduce the issue
3. **Expected Behavior**: What you expected to happen
4. **Actual Behavior**: What actually happened
5. **Environment**:
   - OS and version
   - Go version
   - GoCycled version
6. **Logs/Output**: Any relevant error messages or output

Example:
```
**Bug**: Restore command fails when directory doesn't exist

**Steps to Reproduce**:
1. Trash a file: `rc put /path/to/file.txt`
2. Delete parent directory: `rm -rf /path/to`
3. Try to restore: `rc restore /path/to/file.txt`

**Expected**: Directory should be created and file restored
**Actual**: Error message: "directory does not exist"

**Environment**:
- OS: Ubuntu 22.04
- Go: 1.20
- GoCycled: 1.0.0
```

## Feature Requests

When requesting features:

1. **Use Case**: Describe the problem you're trying to solve
2. **Proposed Solution**: Your idea for solving it
3. **Alternatives**: Other approaches you've considered
4. **Additional Context**: Any other relevant information

## Project Structure

```
GoCycled/
├── cmd/rc/           # Main application entry point
├── pkg/
│   ├── config/       # Configuration management
│   ├── trash/        # Core trash operations
│   └── ui/           # User interface
├── examples/         # Usage examples
├── Makefile          # Build automation
├── README.md         # Project documentation
└── CONTRIBUTING.md   # This file
```

## Code Review Process

1. Maintainers will review your PR
2. Address any feedback or requested changes
3. Once approved, a maintainer will merge your PR
4. Your contribution will be included in the next release

## Questions?

If you have questions:

- Open an issue on GitHub
- Check existing issues and PRs
- Review the documentation in README.md

Thank you for contributing to GoCycled!
