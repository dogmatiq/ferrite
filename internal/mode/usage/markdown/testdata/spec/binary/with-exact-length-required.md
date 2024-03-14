# Environment Variables

## Specification

### `SEED`

> the seed for the random-number generator

The `SEED` variable's value **MUST** be a binary value expressed using the
`base64` encoding scheme with an (unencoded) length of exactly 5 bytes.

```bash
export SEED=wM/6GKc= # (non-normative)
```