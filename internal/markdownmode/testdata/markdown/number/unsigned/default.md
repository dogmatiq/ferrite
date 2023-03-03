# Environment Variables

This document describes the environment variables used by `<app>`.

⚠️ The application may consume other undocumented environment variables; this
document only shows those variables declared using [Ferrite].

Please note that **undefined** variables and **empty strings** are considered
equivalent.

## Index

- [`PPROF_PORT`](#PPROF_PORT) — HTTP port for serving pprof profiling data

## Specification

### `PPROF_PORT`

> HTTP port for serving pprof profiling data

This variable **MAY** be set to a `uint16` value.
If left undefined the default value of `8080` is used.

```bash
export PPROF_PORT=8080 # (default)
```

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
