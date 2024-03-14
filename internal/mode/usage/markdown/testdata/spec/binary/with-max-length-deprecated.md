# Environment Variables

## Specification

### `SEED`

> the seed for the random-number generator

⚠️ The `SEED` variable is **deprecated**; its use is **NOT RECOMMENDED** as it
may be removed in a future version. If defined, the value **MUST** be a binary
value expressed using the `base64` encoding scheme with an (unencoded) length of
10 bytes or fewer.

```bash
export SEED=g5B7W3e8Db5Nbg== # (non-normative)
```
