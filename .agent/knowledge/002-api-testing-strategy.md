---
description: Strategy for backend API testing using k6.
---
# 002: API Testing Strategy

## Status
Proposed

## Context
We need to integrate blackbox E2E API testing into our backend. The requirements are:
- Tests as code
- Can be run locally
- AI friendly
- Functional API testing now, with a path to load testing in the future

Currently, we use:
- **Backend**: Go
- **Frontend**: SvelteKit
- **E2E (Frontend)**: Playwright (TypeScript)

## Decision
We will use **k6** for backend API testing.

## Detailed Analysis

### Option 1: k6 (Recommended)
[k6](https://k6.io/) is a modern load testing tool that is also excellent for functional API testing.

*   **Pros**:
    *   **Dual Purpose**: Scripts written for functional testing can be easily configured for load testing (stress, soak, spike).
    *   **Performance**: Built in Go, extremely performant.
    *   **Language**: Uses JavaScript (ES6), which is familiar to our frontend stack.
    *   **CI/CD**: Designed for automation.
    *   **AI Friendly**: Simple, procedural JS scripts are easy for LLMs to generate and understand.
*   **Cons**:
    *   Introduces a new tool to the stack.
    *   Not a full browser (cannot execute client-side JS), but this is desired for *API* testing.

### Option 2: Playwright
We already use Playwright for frontend E2E. It has a powerful API testing client.

*   **Pros**:
    *   **Unified Stack**: Reuses the existing tool and TypeScript setup.
    *   **Developer Experience**: Excellent assertions and tooling.
*   **Cons**:
    *   **No Load Testing**: Playwright is not designed for load testing. We would need to rewrite tests in another tool (like k6) when load testing becomes necessary.
    *   **Resource Heavy**: Running thousands of Playwright instances is not feasible for load generation.

### Option 3: Go Tests (Standard Lib / Testify)
Native Go testing.

*   **Pros**:
    *   No new languages/tools.
    *   Close to the code.
*   **Cons**:
    *   **Not Blackbox**: Tends to encourage white-box testing or requires complex setup to treat the app as an external black box.
    *   **Load Testing**: Requires writing custom concurrency logic or using a library like `vegeta`, which is less expressive for functional flows than k6.

## Consequences
1.  We will add `k6` to our project.
2.  We will create a `backend/e2e` (or `tests/k6`) directory for these tests.
3.  We will write tests in JavaScript.
4.  We will be ready for load testing immediately.
