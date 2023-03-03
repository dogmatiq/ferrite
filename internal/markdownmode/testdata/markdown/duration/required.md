# Environment Variables

This document describes the environment variables used by `<app>`.

⚠️ Some of the variables have **non-normative** examples. These examples are
syntactically correct but may not be meaningful values for this application.

⚠️ The application may consume other undocumented environment variables; this
document only shows those variables declared using [Ferrite].

Please note that **undefined** variables and **empty strings** are considered
equivalent.

## Index

- [`GRPC_TIMEOUT`](#GRPC_TIMEOUT) — gRPC request timeout

## Specification

### `GRPC_TIMEOUT`

> gRPC request timeout

This variable **MUST** be set to `1ns` or greater.
If left undefined the application will print usage information to `STDERR` then
exit with a non-zero exit code.

```bash
export GRPC_TIMEOUT=640511h56m49.213693952s # (non-normative)
```

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
