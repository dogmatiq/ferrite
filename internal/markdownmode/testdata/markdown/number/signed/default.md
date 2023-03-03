# Environment Variables

This document describes the environment variables used by `<app>`.

⚠️ The application may consume other undocumented environment variables; this
document only shows those variables declared using [Ferrite].

Please note that **undefined** variables and **empty strings** are considered
equivalent.

## Index

- [`WEIGHT`](#WEIGHT) — weighting for this node

## Specification

### `WEIGHT`

> weighting for this node

This variable **MAY** be set to a `int` value.
If left undefined the default value of `+100` is used.

```bash
export WEIGHT=+100 # (default)
```

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
