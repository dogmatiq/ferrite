# Environment Variables

## Specification

### `WEIGHT`

> weighting for this node

The `WEIGHT` variable's value **MUST** be `-10` or greater.

```bash
export WEIGHT=-10 # (non-normative) the minimum accepted value
export WEIGHT=+51 # (non-normative)
export WEIGHT=+72 # (non-normative)
```

<details>
<summary>Signed integer syntax</summary>

Signed integers can only be specified using decimal notation. A leading positive
sign (`+`) is **OPTIONAL**. A leading negative sign (`-`) is **REQUIRED** in
order to specify a negative value.

Internally, the `WEIGHT` variable is represented using a signed 8-bit integer
type (`int8`); any value that overflows this data-type is invalid.

</details>
