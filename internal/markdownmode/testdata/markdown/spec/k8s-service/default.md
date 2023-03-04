# Environment Variables

## Specification

### `REDIS_SERVICE_HOST`

> kubernetes "redis" service host

This variable **MAY** be left undefined, in which case the default value
of `redis.example.org` is used.

```bash
export REDIS_SERVICE_HOST=redis.example.org # (default)
```

### `REDIS_SERVICE_PORT`

> kubernetes "redis" service port

This variable **MAY** be left undefined, in which case the default value
of `6379` is used.

```bash
export REDIS_SERVICE_PORT=6379 # (default)
```
