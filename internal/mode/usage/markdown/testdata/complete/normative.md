# Environment Variables

This document describes the environment variables used by `<app>`.

| Name         | Usage    | Description                                |
| ------------ | -------- | ------------------------------------------ |
| [`READ_DSN`] | required | database connection string for read-models |

> [!TIP]
> If an environment variable is set to an empty value, `<app>` behaves as if
> that variable is left undefined.

## `READ_DSN`

> database connection string for read-models

The `READ_DSN` variable **MUST NOT** be left undefined.

```bash
export READ_DSN=postgres://user:pass@localhost:5432/dbname # user-defined example
```

---

> [!NOTE]
> This document only describes environment variables declared using [Ferrite].
> `<app>` may consume other undocumented environment variables.

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
[`read_dsn`]: #READ_DSN
