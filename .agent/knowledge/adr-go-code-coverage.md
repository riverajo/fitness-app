# ADR: Go Backend Code Coverage

## Status
Proposed

## Context
We need to measure code coverage for the Go backend to ensure test quality and identify untested code paths. The current test verification is manual or basic pass/fail in CI. We need a standardized way to generate, view, and potentially enforce code coverage metrics.

## Options

### Option 1: Standard Go Tooling
Use the built-in `go test` and `go tool cover` commands.
- **Workflow**: 
  - Generate profile: `go test ./... -coverprofile=coverage.out`
  - View summary: `go tool cover -func=coverage.out`
  - View HTML: `go tool cover -html=coverage.out`
- **Pros**: built-in, no dependencies, industry standard.
- **Cons**: CLI output is basic.

### Option 2: `gocov` and `gocov-html`
Use third-party tools to generate more formatted reports.
- **Pros**: Nicer output, generic JSON export.
- **Cons**: Extra dependencies to install/maintain.

### Option 3: SaaS Integration (Codecov/Coveralls)
Upload coverage reports to a SaaS provider.
- **Pros**: Historical tracking, PR comments, nice UI.
- **Cons**: Requires account setup, token management, external service dependency.

### Option 4: Taskfile
Use [Task](https://taskfile.dev/), a modern task runner / build tool that is simpler than Make and uses YAML.
- **Pros**: Easy to read/write (YAML), cross-platform, nice output.
- **Cons**: Environment friction (requires binary installation, path setup).

### Option 5: Makefile (Recommended)
Use a standard `Makefile`.
- **Pros**: Ubiquitous, likely already installed, standard in Go ecosystem.
- **Cons**: Syntax can be finicky (tabs vs spaces), less pretty output than Task.

## Decision
We will use **Option 5: Makefile**.
Due to friction with installing and mounting `task` in the current devcontainer environment, we will revert to the industry-standard `Makefile`. This requires no extra binary installation and should work out-of-the-box.

## Implementation Strategy
1.  Create a `Makefile` in the project root.
2.  Define targets:
    - `test`: Run all tests.
    - `coverage`: Generate coverage profile.
    - `coverage-html`: View HTML report.
