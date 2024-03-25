# Environment Variables

| Name         | Usage    | Description                                | Imported From                      |
| ------------ | -------- | ------------------------------------------ | ---------------------------------- |
| [`READ_DSN`] | required | database connection string for read-models | [Third-party Product](registry:3p) |

## Specification

### `READ_DSN`

> database connection string for read-models

The `READ_DSN` variable **MUST NOT** be left undefined.

```bash
export READ_DSN=foo # (non-normative)
```

This variable is imported from [Third-party Product](registry:3p).

<!-- references -->

[`read_dsn`]: #READ_DSN
[registry:3p]: https://example.org/docs/registry.html
