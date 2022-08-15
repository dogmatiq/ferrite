# Environment Variables

This document describes the environment variables used by `<app>`.

Please note that **undefined** variables and **empty strings** are considered
equivalent.

The application may consume other undocumented environment variables; this
document only shows those variables defined using [Ferrite].

## Index

- [`DEBUG`](#DEBUG) — enable or disable debugging features

## Specification

### `DEBUG`

> enable or disable debugging features

This variable **MAY** be set to one of the values below or left undefined.

```bash
export DEBUG=true
export DEBUG=false
```

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
