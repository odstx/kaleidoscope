# Architecture

## Microservice Proxy

### Overview

The microservice proxy enables the backend to route requests to downstream microservices based on URL path patterns, while maintaining authentication and observability.

### Request Flow

```
Client -> Backend API (/app/:appname/*) -> Microservice Proxy -> Target Service (http://:appname.service/*)
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
