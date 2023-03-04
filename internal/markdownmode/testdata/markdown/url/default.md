# Environment Variables

## Specification

### `API_URL`

> URL of the REST API

This variable **MAY** be set to a non-empty value.
If left undefined the default value is used (see below).

```bash
export API_URL=https://example.org/path # (non-normative) a typical URL for a web page
export API_URL=http://localhost:8080    # (default)
```
