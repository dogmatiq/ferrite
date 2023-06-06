# Environment Variables

## Specification

### `FAVICON`

> the content of the favicon.png file

The `FAVICON` variable **MAY** be left undefined, in which case the default
value of `PGZhdmljb24gY29udGVudD4=` is used. Otherwise, the value **MUST** be a
binary value expressed using the `base64` encoding scheme.

```bash
export FAVICON=PGZhdmljb24gY29udGVudD4= # (default)
```
