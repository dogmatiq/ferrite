# Environment Variables

This document describes the environment variables used by `markdownmode.test`.

The application may consume other undocumented environment variables; this
document only shows those variables defined using [dogmatiq/ferrite].

## Index

- [`DEBUG`](#DEBUG) — enable debug mode

## Specification

### `DEBUG`

> enable or disable debugging features

This variable is **required**, although a default is provided.

- must be one of the values described below

```bash
export DEBUG=true
export DEBUG=false # default value
```

<!-- references -->

[dogmatiq/ferrite]: https://github.com/dogmatiq/ferrite
