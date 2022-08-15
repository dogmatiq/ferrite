# Environment Variables

This document describes the environment variables used by `<app>`.

Please note that **undefined** variables and **empty strings** are considered
equivalent.

The application may consume other undocumented environment variables; this
document only shows those variables defined using [Ferrite].

## Index

- [`PPROF_PORT`](#PPROF_PORT) — HTTP port for serving pprof profiling data

## Specification

### `PPROF_PORT`

> HTTP port for serving pprof profiling data

This variable **MUST** be set to a `uint16` value.
If left undefined the application will print usage information to `STDERR` then
exit with a non-zero exit code.

```bash
export PPROF_PORT=16383 # randomly generated example
```

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
