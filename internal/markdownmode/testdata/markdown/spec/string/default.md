# Environment Variables

## Specification

### `READ_DSN`

> database connection string for read-models

This variable **MAY** be set to a non-empty string.
If left undefined the default value is used (see below).

```bash
export READ_DSN='host=localhost dbname=readmodels user=projector' # (default)
```
