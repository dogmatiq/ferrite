# Environment Variables

## Specification

### `SEED`

> the seed for the random-number generator

⚠️ The `SEED` variable is **deprecated**; its use is **NOT RECOMMENDED** as it
may be removed in a future version. If defined, the value **MUST** have a length
between 5 and 10 bytes.

```bash
export SEED='foo bar' # (non-normative)
```
