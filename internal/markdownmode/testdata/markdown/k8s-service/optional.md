# Environment Variables

This document describes the environment variables used by `<app>`.

⚠️ Some of the variables have **non-normative** examples. These examples are
syntactically correct but may not be meaningful values for this application.

⚠️ The application may consume other undocumented environment variables; this
document only shows those variables declared using [Ferrite].

Please note that **undefined** variables and **empty strings** are considered
equivalent.

## Index

- [`REDIS_SERVICE_HOST`](#REDIS_SERVICE_HOST) — kubernetes "redis" service host
- [`REDIS_SERVICE_PORT`](#REDIS_SERVICE_PORT) — kubernetes "redis" service port

## Specification

### `REDIS_SERVICE_HOST`

> kubernetes "redis" service host

This variable **MAY** be set to a non-empty string or left undefined.

```bash
export REDIS_SERVICE_HOST=foo # (non-normative)
```

### `REDIS_SERVICE_PORT`

> kubernetes "redis" service port

This variable **MAY** be set to a non-empty string or left undefined.

```bash
export REDIS_SERVICE_PORT=foo # (non-normative)
```

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
