# Environment Variables

## Specification

### `API_URL`

> the URL of the REST API

This variable **MAY** be left undefined, in which case the default value
of `http://localhost:8080` is used.
Otherwise, the value **MUST** be a fully-qualified URL.

```bash
export API_URL=https://example.org/path # (non-normative) a typical URL for a web page
export API_URL=http://localhost:8080    # (default)
```
