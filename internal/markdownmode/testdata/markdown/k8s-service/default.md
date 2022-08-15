# Environment Variables

This document describes the environment variables used by `<app>`.

Please note that **undefined** variables and **empty strings** are considered
equivalent.

The application may consume other undocumented environment variables; this
document only shows those variables defined using [Ferrite].

## Index

- [`REDIS_SERVICE_HOST`](#REDIS_SERVICE_HOST) — kubernetes "redis" service host
- [`REDIS_SERVICE_PORT`](#REDIS_SERVICE_PORT) — kubernetes "redis" service port

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

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
