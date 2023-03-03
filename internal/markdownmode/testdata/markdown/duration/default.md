# Environment Variables

This document describes the environment variables used by `<app>`.

⚠️ The application may consume other undocumented environment variables; this
document only shows those variables declared using [Ferrite].

Please note that **undefined** variables and **empty strings** are considered
equivalent.

## Index

- [`GRPC_TIMEOUT`](#GRPC_TIMEOUT) — gRPC request timeout

## Specification

### `GRPC_TIMEOUT`

> gRPC request timeout

This variable **MAY** be set to `1ns` or greater.
If left undefined the default value of `10ms` is used.

Valid durations are a sequence of decimal numbers, each with an optional
fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time
units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

```bash
export GRPC_TIMEOUT=10ms # (default)
```

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
