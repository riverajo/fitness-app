---
description: Strategy for centralized configuration management in the Go backend.
---
# 004: Centralized Configuration Management

## Status
Proposed

## Context
Currently, environment variables are accessed directly via `os.Getenv` throughout the codebase (e.g., in `server.go`, `middleware/auth.go`, `db/mongo.go`, `telemetry/telemetry.go`). This leads to several issues:
- **Dispersed Configuration**: It is difficult to know all the environment variables required to run the application.
- **Late Failures**: Missing environment variables may only be detected when a specific code path is executed, rather than at startup.
- **Type Safety**: Most configuration is treated as strings, requiring manual parsing and validation.
- **Testing**: It is harder to mock configuration when it is hardcoded as `os.Getenv` calls.

We need a system to parse environment variables or CLI arguments and fail fast on startup if the environment is misconfigured.

## Decision
We will implement a centralized **`internal/config`** package and use the **`github.com/caarlos0/env`** library for parsing environment variables into a struct.

## Detailed Analysis

### Option 1: `caarlos0/env` (Selected)
A simple, zero-dependency (other than the library itself) way to parse environment variables into structs using struct tags.

*   **Pros**:
    *   **Simple**: Uses standard Go struct tags.
    *   **Type Safety**: Automatically parses into int, bool, time.Duration, etc.
    *   **Validation**: Supports `required` tag to ensure fail-fast behavior.
    *   **Defaults**: Supports `envDefault` tag for optional values.
    *   **Popularity**: Widely used and well-maintained.
*   **Cons**:
    *   Adds an external dependency (though a lightweight one).

### Option 2: `spf13/viper`
A complete configuration solution for Go applications (JSON, TOML, YAML, ENV, Flags).

*   **Pros**:
    *   Extremely powerful and feature-rich.
    *   Supports live reloading.
*   **Cons**:
    *   **Overkill**: We only need environment variable parsing.
    *   **Complexity**: Larger API surface and dependency footprint.
    *   **Global State**: Often used with a global singleton, which can make testing harder (though not strictly required).

### Option 3: Standard Library (`os.Getenv` + `flag`)
Manually parsing flags and environment variables.

*   **Pros**:
    *   No external dependencies.
*   **Cons**:
    *   **Boilerplate**: Requires writing manual parsing and validation logic for every variable.
    *   **Maintenance**: Harder to maintain as the number of variables grows.

## Implementation Plan
1.  Create a new package `backend/internal/config`.
2.  Define a `Config` struct that holds all configuration values.
    ```go
    type Config struct {
        Port      string `env:"PORT" envDefault:"8080"`
        AppEnv    string `env:"APP_ENV" envDefault:"development"`
        MongoURI  string `env:"MONGO_URI,required"`
        JWTSecret string `env:"JWT_SECRET,required"`
    }
    ```
3.  Implement a `Load()` function that parses the environment variables into this struct using `caarlos0/env`.
4.  Call `config.Load()` at the beginning of `main()`. If it fails, log the error and exit immediately.
5.  Refactor existing components (`db.Connect`, `middleware.AuthMiddleware`, etc.) to accept configuration values as arguments rather than calling `os.Getenv` internally.

## Consequences
1.  **Fail Fast**: The application will refuse to start if required environment variables (like `MONGO_URI` or `JWT_SECRET`) are missing.
2.  **Centralized Definition**: All configuration requirements will be documented in the `Config` struct.
3.  **Dependency Injection**: Components will become more testable as they will rely on passed-in configuration rather than global environment state.
