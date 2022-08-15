# Environment Variables

This document describes the environment variables used by `<app>`.

Please note that **undefined** variables and **empty strings** are considered
equivalent.

The application may consume other undocumented environment variables; this
document only shows those variables defined using [Ferrite].

## Index

- [`WEIGHT`](#WEIGHT) — weighting for this node

## Specification

### `WEIGHT`

> weighting for this node

This variable **MUST** be set to a `int` value.
If left undefined the application will print usage information to `STDERR` then
exit with a non-zero exit code.

```bash
export WEIGHT=-9223372036854775808 # randomly generated example
```

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
