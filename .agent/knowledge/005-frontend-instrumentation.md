# ADR 005: Frontend Instrumentation

## Status
Accepted

## Context
We need to instrument the frontend application to gain visibility into user experience, performance, and errors. The backend is already instrumented with OpenTelemetry. We are using the Grafana stack (Loki, Tempo, Mimir, Grafana) for observability.

## Decision
We will use **Grafana Faro Web SDK** for frontend instrumentation.

## Rationale
1.  **Native Integration**: Faro is designed by Grafana Labs to integrate seamlessly with the Grafana stack (specifically Grafana Agent/Alloy), which we are already using.
2.  **Comprehensive Data**: It collects Web Vitals, logs, exceptions, and custom events out of the box.
3.  **OpenTelemetry Support**: It has built-in support for OpenTelemetry-JS, allowing us to correlate frontend traces with backend traces.
4.  **Ease of Use**: It provides a higher-level abstraction than raw OpenTelemetry JS for RUM (Real User Monitoring) use cases.

## Implementation Details
1.  **Library**: Install `@grafana/faro-web-sdk`.
2.  **Initialization**: Initialize Faro in the SvelteKit app (in `src/routes/+layout.svelte` or a dedicated initialization file).
3.  **Data Transmission**: 
    *   We will configure Faro to send data to a relative URL (e.g., `/faro/collect`).
    *   The Go backend will proxy this route to the actual Grafana Alloy/Agent collector.
    *   This approach avoids CORS issues and simplifies configuration (no need to inject different URLs for different environments into the frontend build).
4.  **Correlation**: Ensure the `traceId` is propagated to backend requests to enable full-stack trace visualization.

## Alternatives Considered
*   **OpenTelemetry JS (Raw)**: More standard, but requires more boilerplate for RUM features (Web Vitals, session tracking) and manual configuration for log shipping to Loki.
*   **Sentry**: Excellent for error tracking, but would require a separate stack/SaaS subscription and wouldn't integrate as "natively" with our self-hosted Grafana stack for logs/traces without additional work.
