# Environment Variables

## Specification

### `WEIGHT`

> weighting for this node

⚠️ The `WEIGHT` variable is **deprecated**; its use is **NOT RECOMMENDED** as it
may be removed in a future version. If defined, the value **MUST** be a whole
number.

```bash
export WEIGHT=-13 # (non-normative)
export WEIGHT=+25 # (non-normative)
```

<details>
<summary>Signed integer syntax</summary>

Signed integers can only be specified using decimal notation. A leading positive
sign (`+`) is **OPTIONAL**. A leading negative sign (`-`) is **REQUIRED** in
order to specify a negative value.

Internally, the `WEIGHT` variable is represented using a signed 8-bit integer
type (`int8`); any value that overflows this data-type is invalid.

</details>
