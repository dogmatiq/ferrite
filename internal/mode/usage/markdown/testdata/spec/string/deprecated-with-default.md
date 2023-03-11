# Environment Variables

## Specification

### `READ_DSN`

> database connection string for read-models

The `READ_DSN` variable **SHOULD** be left undefined, in which case the default
value of `host=localhost dbname=readmodels user=projector` is used.

⚠️ This variable is **deprecated**; its use is discouraged as it may be removed
in a future version.

```bash
export READ_DSN='host=localhost dbname=readmodels user=projector' # (default)
```
