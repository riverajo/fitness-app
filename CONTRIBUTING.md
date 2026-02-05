# Contributing to Fitness App

Thank you for your interest in contributing! This document provides guidelines to ensure a smooth development process.

## Getting Started

1.  **Read the Docs**:
    *   Project Overview: [GEMINI.md](./GEMINI.md)
    *   Backend Setup: [backend/README.md](./backend/README.md)
    *   Frontend Setup: [frontend/README.md](./frontend/README.md)

2.  **Tools**: Ensure you have the following installed:
    *   Go 1.25+
    *   Node.js v20+ & pnpm
    *   Docker & Docker Compose
    *   Lefthook (Git hooks)

## Development Workflow

### 1. Branching Strategy
*   **Main Branch**: `main` (Protected).
*   **Feature Branches**: `feat/description-of-feature`
*   **Fix Branches**: `fix/description-of-bug`
*   **Chore Branches**: `chore/dependency-updates` etc.

### 2. Commits
We follow [Conventional Commits](https://www.conventionalcommits.org/):
*   `feat: add new workout logs`
*   `fix: resolve auth token refresh`
*   `docs: update readme`
*   `refactor: simplify user service`

### 3. Pre-commit Hooks
This project uses **Lefthook**. Hooks will run automatically on commit to ensure:
*   Frontend: Linting (ESLint, Prettier) and Type Checking.
*   Backend: Formatting (`gofmt`), Vet, and Static Analysis.

If a hook fails, fix the issue and commit again.

## Coding Standards

### Backend (Go)
*   **Style**: Follow standard Go idioms. Use `go fmt`.
*   **Architecture**: Respect the separation of concerns:
    *   `graph/`: Resolvers only. No business logic.
    *   `service/`: pure business logic.
    *   `repository/`: Database access only.
*   **Testing**: Write meaningful unit tests. interfaces are used for mocking.

### Frontend (SvelteKit)
*   **Style**: Prettier handles formatting.
*   **State**: Use Svelte 5 Runes (`$state`, `$derived`, `$effect`) for all new local state. Avoid `writable` stores potentially unless for global module state where runes aren't applicable.
*   **Strict Typing**: No `any`. Define interfaces for props and data.

## Pull Requests
1.  Push your branch.
2.  Open a PR against `main`.
3.  Ensure CI passes (Tests, Linting, Build).
4.  Request a review.
