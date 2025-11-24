# ADR 001: Frontend Component Library

## Status
Accepted

## Context
We needed a component library for the Svelte 5 + Tailwind 4 frontend.
Key requirements were:
1.  **Simplicity**: Easy to use drop-in components.
2.  **AI Usability**: Well-known by LLMs.
3.  **Mobile Friendly**: Support for PWA-like features (Bottom Nav, Drawers).
4.  **Tech Stack**: Svelte 5 + Tailwind 4 compatibility.

## Decision
We selected **Flowbite-Svelte**.

## Rationale
- **Mobile Components**: It includes specific components like Bottom Navigation and Drawers which are essential for the "mobile-friendly" requirement.
- **Compatibility**: It supports Tailwind 4 via the `@source` directive and has a path for Svelte 5.
- **AI Support**: Flowbite has excellent documentation structures (`llms.txt`) making it easy for AI agents to generate correct code.

## Consequences
- **Installation**: Requires `flowbite`, `flowbite-svelte`, and `tailwind-variants`.
- **Configuration**: Must use `@source` in CSS for Tailwind v4 to detect classes.
- **Usage**: Future agents should prefer Flowbite components over custom implementations for standard UI elements.
