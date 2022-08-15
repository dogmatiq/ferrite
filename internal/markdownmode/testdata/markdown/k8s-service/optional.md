# Environment Variables

This document describes the environment variables used by `<app>`.

Please note that **undefined** variables and **empty strings** are considered
equivalent.

The application may consume other undocumented environment variables; this
document only shows those variables defined using [Ferrite].

## Index

- [`REDIS_SERVICE_HOST`](#REDIS_SERVICE_HOST) — k8s "redis" service host
- [`REDIS_SERVICE_PORT`](#REDIS_SERVICE_PORT) — k8s "redis" service port

## Specification

### `REDIS_SERVICE_HOST`

> k8s "redis" service host

This variable **MAY** be set to a non-empty string or left undefined.

```bash
export REDIS_SERVICE_HOST=foo # randomly generated example
```

### `REDIS_SERVICE_PORT`

> k8s "redis" service port

This variable **MAY** be set to a non-empty string or left undefined.

```bash
export REDIS_SERVICE_PORT=foo # randomly generated example
```

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
