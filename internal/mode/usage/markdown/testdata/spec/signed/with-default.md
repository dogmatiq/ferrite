# Environment Variables

## `WEIGHT`

> weighting for this node

The `WEIGHT` variable **MAY** be left undefined, in which case the default value
of `+100` is used. Otherwise, the value **MUST** be a whole number.

```bash
export WEIGHT=+100 # (default)
```

<details>
<summary>Signed integer syntax</summary>

Signed integers can only be specified using decimal notation. A leading positive
sign (`+`) is **OPTIONAL**. A leading negative sign (`-`) is **REQUIRED** in
order to specify a negative value.

Internally, the `WEIGHT` variable is represented using a signed 8-bit integer
type (`int8`); any value that overflows this data-type is invalid.

</details>
