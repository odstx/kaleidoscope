# Microservice Registration

## Register a Service Instance

```bash
./kaleidoscope etcd register <app_name> <version> <instance_id> <endpoint>
```

**Example:**
```bash
./kaleidoscope etcd register myapp v1 instance-1 localhost:8080
```

## Deregister a Service Instance

```bash
./kaleidoscope etcd deregister <app_name> <version> <instance_id>
```

**Example:**
```bash
./kaleidoscope etcd deregister myapp v1 instance-1
```

## Options

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--endpoints` | `-e` | `localhost:2379` | etcd server addresses |
| `--username` | `-u` | - | etcd username |
| `--password` | `-p` | - | etcd password |

## etcd Key Structure

Keys are stored in etcd with the following structure:
```
/services/<app_name>/<version>/<instance_id>
```

Value is JSON:
```json
{
  "app_name": "myapp",
  "version": "v1",
  "instance_id": "instance-1",
  "endpoint": "localhost:8080",
  "status": "healthy",
  "metadata": {},
  "registered_at": "2026-04-18T10:00:00Z"
}
```

## Proxy Routing

The proxy routes requests from `/app/<app_name>/<version>/...` to registered service instances.

**Example:**
```
/app/myapp/v1/api/users -> localhost:8080/api/users
```