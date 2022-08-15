# Environment Variables

This document describes the environment variables used by `<app>`.

Please note that **undefined** variables and **empty strings** are considered
equivalent.

The application may consume other undocumented environment variables; this
document only shows those variables defined using [Ferrite].

## Index

- [`READ_DSN`](#READ_DSN) — database connection string for read-models

## Specification

### `READ_DSN`

> database connection string for read-models

This variable **MAY** be set to a non-empty string or left undefined.

```bash
export READ_DSN=foo # randomly generated example
```

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
