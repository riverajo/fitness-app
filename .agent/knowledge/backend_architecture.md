---
description: Overview of backend architecture, patterns, and workflows.
---
# Fitness App Backend Architecture

## Core Pattern: gqlgen Auto-Bind
We use `gqlgen`'s Auto-Bind feature to map GraphQL types directly to our internal domain models. This reduces boilerplate by eliminating the need for separate "Graph" and "Domain" output models.

### 1. Domain Models (`internal/model`)
- **Single Source of Truth**: Core entities (`User`, `WorkoutLog`, `Set`, etc.) are defined in `internal/model`.
- **GraphQL Compatibility**: Struct fields and types must match the GraphQL schema definitions.
    - **IDs**: `ID` fields must be `string` to match GraphQL's `ID` scalar.
    - **Tags**: Use `json` tags to match GraphQL field names.
- **Exclusions**: Input types (e.g., `CreateWorkoutLogInput`) are **not** defined here to avoid circular dependencies. We use the generated types from `graph/model`.

### 2. GraphQL Layer (`graph/`)
- **Resolvers**:
    - **Inputs**: Accept generated Input types from `graph/model` (e.g., `model1.CreateWorkoutLogInput`).
    - **Outputs**: Return internal Domain models directly (e.g., `*internalModel.WorkoutLog`).
    - **Mapping**:
        - **Output**: Automatic. No manual mapping required.
        - **Input**: Manual. You must map the generated Input structs to your Domain structs within the resolver (or a helper).

### 3. Repository Layer (`internal/repository`)
- **Responsibility**: Handles the impedance mismatch between the Application Layer (String IDs) and Database Layer (ObjectID).
- **ID Conversion**:
    - **Reads**: When fetching from Mongo, convert the `_id` (ObjectID) to the `ID` (string) field in the domain model.
    - **Writes**: When saving to Mongo, convert the `ID` (string) to `_id` (ObjectID).
    - **Tips**: Use `primitive.ObjectIDFromHex(id)` for string->ObjectID and `objID.Hex()` for ObjectID->string.

### 4. Authentication
- **Context**: Authenticated User ID is available in the context via `middleware.UserIDKey`.
- **Type**: The User ID in context is a `string`.

### 5. Workflow for New Features
1.  **Schema**: Define types in `graph/schema.graphqls`.
2.  **Model**: Create/Update the struct in `internal/model` to match.
3.  **Generate**: Run `go run github.com/99designs/gqlgen generate`.
4.  **Repository**: Implement DB logic, handling ID conversion.
5.  **Resolver**: Implement resolver using the Internal Model directly.
