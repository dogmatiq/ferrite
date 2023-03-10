# Environment Variables

## Specification

### `DEBUG`

> enable or disable debugging features

The `DEBUG` variable **SHOULD** be left undefined, in which case the default
value of `false` is used. Otherwise, the value **MUST** be either `true` or
`false`.

⚠️ This variable is **deprecated**; its use is discouraged as it may be removed
in a future version.

```bash
export DEBUG=true
export DEBUG=false # (default)
```
