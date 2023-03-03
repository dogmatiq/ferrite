# Environment Variables

This document describes the environment variables used by `<app>`.

⚠️ The application may consume other undocumented environment variables; this
document only shows those variables declared using [Ferrite].

Please note that **undefined** variables and **empty strings** are considered
equivalent.

## Index

- [`READ_DSN`](#READ_DSN) — database connection string for read-models

## Specification

### `READ_DSN`

> database connection string for read-models

This variable **MAY** be set to a non-empty string.
If left undefined the default value is used (see below).

```bash
export READ_DSN='host=localhost dbname=readmodels user=projector' # (default)
```

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
