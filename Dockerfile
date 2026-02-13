# Stage 1: Build Frontend
FROM node:24-alpine AS frontend-builder
WORKDIR /app/frontend
# Install pnpm
RUN npm install -g pnpm
COPY frontend/package.json frontend/pnpm-lock.yaml frontend/svelte.config.js ./
RUN pnpm install --frozen-lockfile
COPY frontend/ .
RUN pnpm run build

# Stage 2: Build Backend
FROM golang:1.26-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
# Copy frontend build artifacts to backend/public for embedding
COPY --from=frontend-builder /app/frontend/build ./public
# Build the binary
# CGO_ENABLED=0 for static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o fitness-app .

# Stage 3: Final Image
FROM alpine:latest
WORKDIR /app
# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates
# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
# Copy the binary from the backend builder
COPY --from=backend-builder /app/backend/fitness-app .
# Expose the port (assuming 8080, adjust if needed)
EXPOSE 8080
CMD ["./fitness-app"]
