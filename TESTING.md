# Testing Strategy

This repository employs a multi-layered testing strategy to ensure reliability across the stack.

## Backend (Go)

We use the standard `testing` package along with `testify` for assertions and mocking.

### 1. Unit Tests
Located alongside the code they test (e.g., `service_test.go`).
*   **Goal**: Verify business logic in isolation.
*   **Mocks**: Use interface-based mocking (often generated or manually defined).
*   **Command**:
    ```bash
    go test ./...
    ```

### 2. Integration Tests
Tests that talk to a real (dockerized) database.
*   **Setup**: Uses `testcontainers-go` to spin up a MongoDB instance per test suite.
*   **Command**: Same as unit tests, but may take longer.
    ```bash
    make test
    ```

### 3. Coverage
To view test coverage visually:
```bash
make coverage-html
```

### 4. Load Tests (k6)
Load testing for API endpoints (e.g., login, register) is handled by [k6](https://k6.io/).
*   **Location**: `backend/k6/`
*   **Prerequisite**: Install k6 (see [docs](https://k6.io/docs/get-started/installation/)). The backend must be running (`docker compose up api`).
*   **Command**:
    ```bash
    # Run from repository root
    k6 run backend/k6/main.js
    ```

## Frontend (SvelteKit)

### 1. Unit Tests (Vitest)
Tests for individual components and utility functions.
*   **Goal**: Ensure isolated components render and behave correctly.
*   **Command**:
    ```bash
    cd frontend
    pnpm run test:unit
    ```

### 2. End-to-End Tests (Playwright)
Simulates real user interactions across the full stack.
*   **Prerequisite**: The backend and frontend servers must be running (or handled by the test runner if configured).
*   **Goal**: critical user flows (Login, Register, Create Workout).
*   **Command**:
    ```bash
    cd frontend
    # Runs tests in headless mode
    pnpm run test:e2e
    
    # Opens the Playwright UI for debugging
    pnpm run test:e2e -- --ui
    ```

## CI/CD Pipeline

On every Pull Request, the following checks run automatically:
1.  **Backend**: `go test`, `go vet`, `staticcheck`.
2.  **Frontend**: `eslint`, `svelte-check`, `vitest`.
3.  **E2E**: `playwright` (Ensures no regression in critical flows).
