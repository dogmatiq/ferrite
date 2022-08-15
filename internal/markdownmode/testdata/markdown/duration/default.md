# Environment Variables

This document describes the environment variables used by `<app>`.

Please note that **undefined** variables and **empty strings** are considered
equivalent.

The application may consume other undocumented environment variables; this
document only shows those variables defined using [Ferrite].

## Index

- [`GRPC_TIMEOUT`](#GRPC_TIMEOUT) — gRPC request timeout

## Specification

### `GRPC_TIMEOUT`

> gRPC request timeout

This variable **MAY** be set to `1ns` or greater.
If left undefined the default value of `10ms` is used.

```bash
export GRPC_TIMEOUT=10ms # (default)
```

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
