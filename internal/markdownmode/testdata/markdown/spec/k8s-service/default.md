# Environment Variables

## Specification

### `REDIS_SERVICE_HOST`

> kubernetes "redis" service host

This variable **MAY** be set to a non-empty string.
If left undefined the default value is used (see below).

```bash
export REDIS_SERVICE_HOST=redis.example.org # (default)
```

### `REDIS_SERVICE_PORT`

> kubernetes "redis" service port

This variable **MAY** be set to a non-empty string.
If left undefined the default value is used (see below).

```bash
export REDIS_SERVICE_PORT=6379 # (default)
```
