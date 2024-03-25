# Environment Variables

## `API_URL`

> the URL of the REST API

⚠️ The `API_URL` variable is **deprecated**; its use is **NOT RECOMMENDED** as
it may be removed in a future version. If defined, the value **MUST** be a
fully-qualified URL.

```bash
export API_URL=https://example.org/path # (non-normative) a typical URL for a web page
```

<details>
<summary>URL syntax</summary>

A fully-qualified URL includes both a scheme (protocol) and a hostname. URLs are
not necessarily web addresses; `https://example.org` and
`mailto:contact@example.org` are both examples of fully-qualified URLs.

</details>
