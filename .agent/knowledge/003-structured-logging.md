---
description: Strategy for structured logging in the Go backend.
---
# 003: Structured Logging Strategy

## Status
Proposed

## Context
We need to implement structured logging in our Go backend to improve observability and debuggability. The key criteria for selecting a logging solution are:
- **Simplicity**: Minimal configuration and dependencies.
- **Ease of Use**: Intuitive API for developers.
- **JSON Structure**: Native support for JSON output for machine parsing.
- **AI Aptitude**: The library and standard should be well-understood by AI models (common patterns, standard libraries).
- **Readability**: Logs should be readable by both AI agents and humans.

## Decision
We will use the standard library **`log/slog`** (introduced in Go 1.21).

## Detailed Analysis

### Option 1: `log/slog` (Recommended)
Go 1.21 introduced structured logging to the standard library.

*   **Pros**:
    *   **Standard**: It is the official Go standard. This ensures long-term support and compatibility.
    *   **Simplicity**: Zero external dependencies.
    *   **AI Aptitude**: As the standard library, it is the default "correct" way to log in Go moving forward, making it highly predictable for AI.
    *   **Performance**: Highly optimized.
    *   **Structure**: Built-in `JSONHandler` meets the JSON requirement out of the box.
*   **Cons**:
    *   Newer than established libraries like Zap or Zerolog (though rapidly becoming the default).

### Option 2: Uber Zap
A popular, high-performance structured logging library.

*   **Pros**:
    *   Extremely fast.
    *   Type-safe field construction.
*   **Cons**:
    *   **Complexity**: API can be verbose (`zap.String("key", "val")`) or requires a "Sugared" logger for simplicity, which splits the API surface.
    *   External dependency.

### Option 3: Zerolog
A zero-allocation JSON logger.

*   **Pros**:
    *   Great developer experience (chained API).
    *   Fast.
*   **Cons**:
    *   External dependency.
    *   Non-standard API style (chained methods vs variadic args).

## Consequences
1.  We will configure the global logger in `main.go` to use `slog.New(slog.NewJSONHandler(os.Stdout, nil))`.
2.  We will replace existing `log.Println` and `fmt.Println` calls with `slog.Info`, `slog.Error`, etc.
3.  We will use structured attributes (e.g., `slog.String("user_id", id)`) instead of formatted strings where appropriate.
