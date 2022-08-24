# Environment Variables

This document describes the environment variables used by `<app>`.

Please note that **undefined** variables and **empty strings** are considered
equivalent.

The application may consume other undocumented environment variables; this
document only shows those variables defined using [Ferrite].

## Index

- [`PORT`](#PORT) — listen port for the HTTP server

## Specification

### `PORT`

> listen port for the HTTP server

This variable **MUST** be set to a non-empty string.
If left undefined the application will print usage information to `STDERR` then
exit with a non-zero exit code.

```bash
export PORT=foo # randomly generated example
```

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
