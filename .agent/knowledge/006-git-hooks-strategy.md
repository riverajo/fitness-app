# ADR 006: Git Hook Strategy

## Status
Accepted

## Context
We need to enforce code quality standards (linting, formatting, testing) automatically before changes are committed to the repository. This "shift-left" approach ensures that bad code (syntax errors, failing tests, unformatted files) never enters version control, saving CI resources and developer time.

The repository is a monorepo-style structure with:
- `backend/` (Go)
- `frontend/` (SvelteKit/TypeScript)
- Root-level configuration (Docker, etc.)

## Options

### 1. Lefthook
A fast, polyglot git hook manager written in Go.
- **Pros**: 
    - Extremely fast (parallel execution).
    - Single binary, easy to install (Go, Node, Ruby, Homebrew).
    - Excellent support for monorepos (can run commands in specific subdirectories).
    - Configuration in a single `lefthook.yml` at the root.
- **Cons**: 
    - Another tool to install (though can be a devDependency in frontend).

### 2. Husky
Review standard for Node.js projects.
- **Pros**: 
    - Ubiquitous in the JS ecosystem.
    - Integration with `package.json`.
- **Cons**: 
    - Slower than Lefthook.
    - Root of the repo is not a Node project (frontend is a subdir), making setup slightly awkward or requiring a root `package.json`.
    - "Lint-staged" setup can be complex for polyglot repos.

### 3. pre-commit
Python-based framework.
- **Pros**: 
    - Large library of existing hooks.
    - Isolated environments for hooks.
- **Cons**: 
    - Requires Python.
    - Can be slow (environment creation).
    - Overkill if we just want to run local tools (golangci-lint, eslint) we already found.

### 4. Native Git Hooks (Shell Scripts)
- **Pros**: No dependencies.
- **Cons**: Hard to manage, share, and maintain. No parallelization constraints.

## Decision
We will use **Lefthook**.

### Rationale
- **Performance**: Lefthook is significantly faster, which is critical for developer experience.
- **Monorepo Support**: It natively handles our `backend/` vs `frontend/` split gracefully without needing a root `package.json`.
- **Ecosystem Fit**: Being a Go binary, it fits well with our Backend stack, but installs easily via `npm` for Frontend workflows if desired.

## Implementation Details
- Install `lefthook` (via `go install` or `programs`).
- Create `lefthook.yml` at project root.
- Configure `pre-commit` to run:
    - **Backend**: `golangci-lint run`, `go test ./...`
    - **Frontend**: `npm run check` (permissions allowing), `eslint`.
- We will configure it to generally only check staged files where possible, or run fast checks on the whole project if speed permits.
