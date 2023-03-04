# Environment Variables

This document describes the environment variables used by `<app>`.

If any of the environment variable values do not meet the requirements herein,
the application will print usage information to `STDERR` then exit with a
non-zero exit code. Please note that **undefined** variables and **empty**
values are considered equivalent.

⚠️ This document includes **non-normative** example values. While these values
are syntactically correct, they may not be meaningful to this application.

⚠️ The application may consume other undocumented environment variables; this
document only shows those variables declared using [Ferrite].

The key words **MUST**, **MUST NOT**, **REQUIRED**, **SHALL**, **SHALL NOT**,
**SHOULD**, **SHOULD NOT**, **RECOMMENDED**, **MAY**, and **OPTIONAL** in this
document are to be interpreted as described in [RFC 2119].

## Index

- [`READ_DSN`](#READ_DSN) — database connection string for read-models

## Specification

### `READ_DSN`

> database connection string for read-models

This variable **MUST** be set to a non-empty string.

```bash
export READ_DSN=foo # (non-normative)
```

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
[rfc 2119]: https://www.rfc-editor.org/rfc/rfc2119.html
