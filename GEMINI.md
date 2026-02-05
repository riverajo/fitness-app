# GEMINI - Project Context & Memory

This file serves as the "External Brain" for the project. It documents the persistent context, architectural decisions, and development guidelines.

## System Overview

**Tech Stack**
- **Backend**: Go 1.25+, GraphQL (`99designs/gqlgen`), MongoDB (Official Driver), `golang-jwt`.
- **Frontend**: SvelteKit (Svelte 5), TypeScript, TailwindCSS v4, Flowbite Svelte, URQL (GraphQL Client).
- **Infrastructure**: Docker Compose (Dev), OpenTelemetry (Observability).
- **Tooling**: Lefthook (Git Hooks), Renovate (Deps).

**Core Domain**: Fitness tracking application focusing on workouts, exercises, and user progress.

## Architecture Map

### Data Flow
1.  **UI Layer**: User interacts with Svelte components. State is manged via Svelte 5 runes and persistent `authStore`.
2.  **Data Fetching**: `URQL` client sends GraphQL queries/mutations to `/query`. Authentication is handled via `authExchange` (JWT + Refresh Token).
3.  **API Gateway**: Go server (`backend/server.go`) receives requests.
    *   Global Middleware: Logging, Tracing, Auth (Context Extraction).
4.  **GraphQL Layer**: `backend/graph` resolvers handle specific operations.
5.  **Service Layer**: `backend/internal/service` contains the business logic.
6.  **Repository Layer**: `backend/internal/repository` handles database I/O.
7.  **Database**: MongoDB stores data as documents (Workouts, Users, Exercises).

## The 'Why' - Design Decisions

- **Dual-Token Auth**: Used to balance security and convenience. Access tokens (short-lived) are used for API calls, while HttpOnly cookies (refresh tokens) allow for silent persistence without exposing long-lived credentials to XSS.
- **GraphQL**: Enforces a strict contract between the Go backend and TypeScript frontend. Code generation ensures types are always in sync.
- **Go Backend Structure**: Follows Standard Go Project Layout (roughly).
    - `internal/`: Private application code to prevent external import.
    - `graph/`: strictly for GraphQL schema and wiring.
    - `service/` vs `repository/`: Strict separation of concerns. Services do not know about DB implementations; Repositories do not handle business rules.
- **Runes (Svelte 5)**: We use the new Svelte 5 reactivity system (runes) for granular state management inside `frontend/src/state`.

## Current Trajectory

**Recent Features**
- Migration to Dual-Token Authentication.
- Implementation of strict strict `WeightUnit` Enum handling across full stack.
- Frontend Refactors (Runes adoption).

**Next Steps & Logical Evolution**
- **Feature Completion**: Fill out CRUD for Workouts (Filtering, History).
- **Testing**: Expand E2E coverage for complex flows (Auth expiry, Offline mode).
- **Observability**: Tune OpenTelemetry traces for better debugging.

## Agent Instructions

**Coding Standards**
- **Go**: Interface-based dependency injection. Always define interfaces in the consumer package.
- **Svelte**: Use Runes (`$state`, `$derived`) over stores where possible. Use `lang="ts"`.
- **GraphQL**: If changing schema, run `go run github.com/99designs/gqlgen generate` (or project equivalent task) AND frontend codegen.

**Testing Requirements**
- **Frontend**: Run `npm run test` (Unit + E2E) before finishing tasks.
- **Backend**: Run `go test ./...` in the `backend` directory.

**Workflow**
- When working on the frontend, check `frontend/package.json` scripts.
- When working on the backend, verify changes with `Makefile`.
