# Environment Variables

## Specification

### `API_URL`

> the URL of the REST API

The `API_URL` variable **SHOULD** be left undefined, in which case the default
value of `http://localhost:8080` is used. Otherwise, the value **MUST** be a
fully-qualified URL.

⚠️ This variable is **deprecated**; its use is discouraged as it may be removed
in a future version.

```bash
export API_URL=http://localhost:8080    # (default)
export API_URL=https://example.org/path # (non-normative) a typical URL for a web page
```

<details>
<summary>URL syntax</summary>

A fully-qualified URL includes both a scheme (protocol) and a hostname. URLs are
not necessarily web addresses; `https://example.org` and
`mailto:contact@example.org` are both examples of fully-qualified URLs.

</details>
