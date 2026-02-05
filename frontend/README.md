# Fitness App Frontend

A modern fitness tracking application built with **SvelteKit (Svelte 5)**, **TypeScript**, and **TailwindCSS**.

## Tech Stack

- **Framework**: [SvelteKit](https://kit.svelte.dev/) (Svelte 5 with Runes)
- **Styling**: [TailwindCSS v4](https://tailwindcss.com/) & [Flowbite Svelte](https://flowbite-svelte.com/)
- **Data Fetching**: [URQL](https://formidable.com/open-source/urql/) (GraphQL)
- **State Management**: Svelte 5 Runes (`$state`, `$derived`)
- **Testing**: [Playwright](https://playwright.dev/) (E2E), [Vitest](https://vitest.dev/) (Unit)

## Getting Started

### Prerequisites

- Node.js (v20+)
- pnpm (Recommended) or npm

### Installation

1. Install dependencies:

   ```bash
   pnpm install
   ```

2. Generate GraphQL types (ensure backend is running or schema is available):
   ```bash
   pnpm run codegen
   ```

### Development

Start the development server:

```bash
pnpm run dev
```

The application will be available at `http://localhost:5173`.

## Testing

**Unit Tests**:

```bash
pnpm run test:unit
```

**E2E Tests**:

```bash
# Ensure backend is running!
pnpm run test:e2e
```

**Run All Tests**:

```bash
pnpm run test
```

## Key Directories

- `src/lib`: Shared utilities and UI components.
- `src/routes`: SvelteKit File-system routing.
- `src/state`: Global state management using Svelte 5 Runes.
- `src/lib/gql`: Generated GraphQL types and operations.
