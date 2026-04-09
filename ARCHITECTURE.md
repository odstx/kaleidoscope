# Architecture

## Microservice Proxy

### Overview

The microservice proxy enables the backend to route requests to downstream microservices based on URL path patterns, while maintaining authentication and observability.

### Request Flow

```mermaid
sequenceDiagram
    participant Client
    participant Backend
    participant Proxy as Microservice Proxy
    participant Database
    participant Service as :appname.service

    Client->>Backend: GET /app/myapp/api/users
    Backend->>Proxy: Pass request
    Note over Proxy: Authentication Check
    alt JWT or Hawk
        Client->>Backend: Bearer token / Hawk auth
    else No Auth
        Backend-->>Client: 401 Unauthorized
    end
    Note over Proxy: Query User Info
    Proxy->>Database: SELECT username WHERE uid = ?
    Database-->>Proxy: username
    Note over Proxy: Create OTEL Span
    Note over Proxy: Forward Request
    Proxy->>Service: GET /api/users (X-UID, X-Username, trace headers)
    Service-->>Proxy: Response
    Proxy-->>Client: Forward Response
```

### Routing Rules

| Path Pattern | Target URL | Example |
|--------------|------------|---------|
| `/app/myapp/` | `http://myapp.service/` | - |
| `/app/myapp/api/users` | `http://myapp.service/api/users` | - |

### Configuration

In `config.yaml`:

```yaml
microservice:
  enabled: true
  service_domain: "service"  # Kubernetes service domain suffix
```

### Authentication

- All `/app/*` routes require user authentication
- Supports JWT (Bearer token) and Hawk authentication
- Authentication is enforced via `CombinedAuth` middleware

### Headers Forwarding

The proxy forwards the following headers to downstream services:

| Header | Description |
|--------|-------------|
| `X-UID` | User ID from database |
| `X-Username` | Username from database |
| `W3C Trace Context` | OTEL propagation headers |

### Observability

- Creates OTEL span with attributes: `app.name`, `target.url`, `user.id`
- Injects trace context into proxied requests for distributed tracing

### Security

- Path traversal prevention via `path.Clean()`
- Authentication required for all proxied requests
- User identity passed via headers (not in URL)

## Frontend Micro-Frontend Architecture

### Overview

The frontend uses a Web Component-based micro-frontend architecture. Pod is built as a custom element and loaded dynamically by the MicroAppPage.

### Request Flow

```mermaid
sequenceDiagram
    participant User
    participant Browser
    participant Frontend as React App
    participant Backend as Kaleidoscope API
    participant Pod as Pod Web Component

    User->>Frontend: Navigate to /app/pod
    Note over Frontend: Check authentication
    alt Not authenticated
        Frontend->>User: Redirect to Login
    else Authenticated
        Frontend->>Browser: Load /app/pod.js (dynamic script)
        Browser->>Frontend: Script loaded
        Frontend->>Browser: customElements.whenDefined('pod-app')
        Frontend->>Browser: Create pod-app element
        Browser->>Pod: Mount React app in shadow DOM
    end
```

### Build Output

```mermaid
graph LR
    A[pod/src/pod-wc.tsx] --> B[Vite Build]
    B --> C[frontend/dist/app/pod.js]
    C --> D[Browser loads as ES module]
    D --> E[customElements.define]
```

### Configuration

- Build output: `frontend/dist/app/:appname.js`
- Custom element tag: `:appname-app`
- Route pattern: `/app/:appname`
