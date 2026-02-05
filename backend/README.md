# Fitness App Backend

The GraphQL API gateway and business logic server for the Fitness App.

## Tech Stack

- **Language**: Go 1.25+
- **API**: GraphQL ([99designs/gqlgen](https://github.com/99designs/gqlgen))
- **Database**: MongoDB (Official Go Driver)
- **Auth**: Dual-token (JWT Access Token + HttpOnly Refresh Cookie)
- **Observability**: OpenTelemetry

## Getting Started

### Prerequisites

- Go 1.25 or later
- Docker (for MongoDB)
- Air (optional, for hot reload)

### Development

1. **Start Infrastructure**:
   Ensure MongoDB is running (usually via `docker compose up -d mongo` from root).

2. **Run Server**:
   ```bash
   # Standard run
   go run server.go

   # With hot reload (recommended)
   air
   ```
   The server listens on port `8080`.

3. **GraphQL Playground**:
   Visit `http://localhost:8080/` to access the GraphiQL playground.

## Code Generation

When modifying `graph/schema.graphqls`, regenerate the Go resolvers:

```bash
go run github.com/99designs/gqlgen generate
```

## Testing

Run unit and integration tests:

```bash
# Run all tests
make test

# Generate coverage report
make coverage-html
```

## Directory Structure

- `server.go`: Application entrypoint.
- `graph/`: GraphQL schema, resolvers, and generated code.
- `internal/`: Private application code.
    - `api/`: HTTP handlers and middleware.
    - `service/`: Business logic.
    - `repository/`: Database interactions.
    - `model/`: Domain models.
