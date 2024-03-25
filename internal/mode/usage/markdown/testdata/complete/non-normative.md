# Environment Variables

This document describes the environment variables used by `<app>`.

| Name         | Usage    | Description                                |
| ------------ | -------- | ------------------------------------------ |
| [`READ_DSN`] | required | database connection string for read-models |

> [!WARNING]
> This document only shows environment variables declared using [Ferrite].
> `<app>` may consume other undocumented environment variables.

## Specification

All environment variables described below must meet the stated requirements.
Otherwise, `<app>` prints usage information to `STDERR` then exits.
**Undefined** variables and **empty** values are equivalent.

⚠️ This section includes **non-normative** example values. These examples are
syntactically valid, but may not be meaningful to `<app>`.

The key words **MUST**, **MUST NOT**, **REQUIRED**, **SHALL**, **SHALL NOT**,
**SHOULD**, **SHOULD NOT**, **RECOMMENDED**, **MAY**, and **OPTIONAL** in this
document are to be interpreted as described in [RFC 2119].

### `READ_DSN`

> database connection string for read-models

The `READ_DSN` variable **MUST NOT** be left undefined.

```bash
export READ_DSN=foo # (non-normative)
```

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
[`read_dsn`]: #READ_DSN
[rfc 2119]: https://www.rfc-editor.org/rfc/rfc2119.html
