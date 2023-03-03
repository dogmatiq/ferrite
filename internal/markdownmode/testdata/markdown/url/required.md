# Environment Variables

This document describes the environment variables used by `<app>`.

⚠️ Some of the variables have **non-normative** examples. These examples are
syntactically correct but may not be meaningful values for this application.

⚠️ The application may consume other undocumented environment variables; this
document only shows those variables declared using [Ferrite].

Please note that **undefined** variables and **empty strings** are considered
equivalent.

## Index

- [`API_URL`](#API_URL) — URL of the REST API

## Specification

### `API_URL`

> URL of the REST API

This variable **MUST** be set to a non-empty value.
If left undefined the application will print usage information to `STDERR` then
exit with a non-zero exit code.

```bash
export API_URL=https://example.org/path # (non-normative) a typical URL for a web page
```

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite