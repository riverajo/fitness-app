# Backend Architecture & Feature Guide

This document deep divides into the backend structure to help you understand where code lives and how to add new features.

## Layered Architecture

We follow a strict layered architecture to enforce separation of concerns:

1.  **Transport Layer (GraphQL)**: `backend/graph`
    *   **Responsibility**: Parsing requests, input validation, calls service layer.
    *   **Rules**: No business logic here. Only conversion between GraphQL types and Domain models.

2.  **Service Layer**: `backend/internal/service`
    *   **Responsibility**: Core business logic, authorization checks, transaction management.
    *   **Rules**: "Pure" Go code. Should not depend on specific database implementations (uses Repository interfaces).

3.  **Repository Layer**: `backend/internal/repository`
    *   **Responsibility**: CRUD operations, database queries.
    *   **Rules**: No business logic. Just data access.

4.  **Domain Models**: `backend/internal/model`
    *   **Responsibility**: Structs and Types shared across the application.

## Request Lifecycle

1.  **Request**: Inbound HTTP POST to `/query`.
2.  **Middleware**: `internal/api/middleware` extracts JWT, sets user context, starts tracing.
3.  **Resolver**: `graph/schema.resolvers.go` receives the request.
4.  **Service**: Resolver calls `Service.CreateX()`.
5.  **Repository**: Service calls `Repository.InsertX()`.
6.  **Database**: Data persisted in MongoDB.

## How to Add a New Feature

**Example**: Adding a "Goal" feature.

1.  **Define Schema**:
    *   Edit `graph/schema.graphqls`.
    *   Add `type Goal`, `input CreateGoalInput`, and mutations.

2.  **Generate Code**:
    *   Run `go run github.com/99designs/gqlgen generate`.

3.  **Define Model**:
    *   Create `internal/model/goal.go`.

4.  **Create Repository**:
    *   Define interface: `internal/repository/goal_repository.go`.
    *   Implement Mongo version: `internal/repository/mongo_goal_repository.go`.

5.  **Create Service**:
    *   Define interface: `internal/service/goal_service.go`.
    *   Implement logic: `internal/service/goal_service_impl.go` (injects GoalRepository).

6.  **Wire It Up**:
    *   In `server.go`, initialize the new Repository and Service.
    *   Inject the Service into the `Resolver` struct.

7.  **Implement Resolver**:
    *   In `graph/schema.resolvers.go`, update the generated methods to call your new Service.
