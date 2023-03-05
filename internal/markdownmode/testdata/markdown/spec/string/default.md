# Environment Variables

## Specification

### `READ_DSN`

> database connection string for read-models

The `READ_DSN` variable **MAY** be left undefined, in which case the default
value of `host=localhost dbname=readmodels user=projector` is used.

```bash
export READ_DSN='host=localhost dbname=readmodels user=projector' # (default)
```
