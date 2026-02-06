# Fitness App - Agentic AI Playground

## Introduction
The primary goal of this project is to explore the world of **agentic AI coding**. It serves as a sandbox for experimenting with AI-driven development workflows, allowing for hands-on experience with how AI agents can assist in building, refactoring, and maintaining a full-stack application.

Most of the technology stack was selected specifically for **play and learning**, rather than purely for enterprise constraints. This allows for exploring modern, interesting, and sometimes bleeding-edge tools to see how they interact with AI coding assistants.

## Technology Stack

### Backend
The backend is built with **Go**, focusing on performance and type safety.
*   **Language:** Go
*   **API:** GraphQL using [gqlgen](https://github.com/99designs/gqlgen)
*   **Database:** MongoDB (accessed via the official Go driver)
*   **Authentication:** JWT-based auth
*   **Testing:** Testcontainers for robust integration testing
*   **API Testing:** k6 for performance and load testing

### Frontend
The frontend is a modern, reactive web application built with **Svelte**.
*   **Framework:** Svelte 5 & SvelteKit
*   **Styling:** Tailwind CSS & Flowbite
*   **GraphQL Client:** urql
*   **Testing:** Playwright (E2E) & Vitest (Unit)

## Getting Started

The entire development stack is containerized using Docker Compose for easy setup.

### Prerequisites
*   Docker & Docker Compose
*   [k6](https://k6.io/) (optional, for load testing)

### Running the Application
To start the full stack (Database, Backend API, Frontend, and Adminer):

```bash
docker compose up --build
```

This will spin up the following services:
*   **Frontend:** [http://localhost:5173](http://localhost:5173)
*   **Backend API:** [http://localhost:8080](http://localhost:8080) (GraphQL Playground available in dev mode)
*   **MongoDB:** Exposed on port `27017`

### Development
*   **Backend:** The backend service uses `air` for hot-reloading. Changes to `.go` files will automatically trigger a rebuild.
*   **Frontend:** The frontend container mounts the source code, enabling hot module replacement (HMR) via Vite.

## Documentation

*   [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution guidelines
*   [TESTING.md](TESTING.md) - Testing strategy and commands
*   [GEMINI.md](GEMINI.md) - Project context and memory
*   [backend/ARCHITECTURE.md](backend/ARCHITECTURE.md) - Backend architectural details
*   [backend/README.md](backend/README.md) - Backend specific documentation
*   [frontend/README.md](frontend/README.md) - Frontend specific documentation
