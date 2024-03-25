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

The `READ_DSN` variable is **required**. It must have a non-empty value.

```bash
export READ_DSN=foo # (non-normative)
```

---

> [!NOTE]
> This document only describes environment variables declared using [Ferrite].
> `<app>` may consume other undocumented environment variables.

> [!IMPORTANT]
> Some of the example values given in this document are **non-normative**.
> Although these values are syntactically valid, they may not be meaningful to
> `<app>`.

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
[`read_dsn`]: #READ_DSN
