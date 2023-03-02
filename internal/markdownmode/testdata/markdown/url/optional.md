# Environment Variables

This document describes the environment variables used by `<app>`.

Please note that **undefined** variables and **empty strings** are considered
equivalent.

The application may consume other undocumented environment variables; this
document only shows those variables defined using [Ferrite].

## Index

- [`API_URL`](#API_URL) — URL of the REST API

## Specification

### `API_URL`

> URL of the REST API

This variable **MAY** be set to a non-empty value or left undefined.

```bash
export API_URL=https://example.org/path # randomly generated example
```

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
