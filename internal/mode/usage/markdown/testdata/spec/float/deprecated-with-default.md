# Environment Variables

## Specification

### `WEIGHT`

> weighting for this node

The `WEIGHT` variable **SHOULD** be left undefined, in which case the default
value of `+123.5` is used. Otherwise, the value **MUST** be a number with an
**OPTIONAL** fractional part.

⚠️ This variable is **deprecated**; its use is discouraged as it may be removed
in a future version.

```bash
export WEIGHT=+123.5 # (default)
```

<details>
<summary>Floating-point syntax</summary>

Floating-point values can be specified using decimal (base-10) or hexadecimal
(base-16) notation, and may use scientific notation. A leading positive sign
(`+`) is **OPTIONAL**. A leading negative sign (`-`) is **REQUIRED** in order to
specify a negative value.

Internally, the `WEIGHT` variable is represented using a 32-bit floating point
type (`float32`); any value that overflows this data-type is invalid. Values are
rounded to the nearest floating-point number using IEEE 754 unbiased rounding.

The non-finite values `NaN`, `+Inf` and `-Inf` are not accepted.

</details>
