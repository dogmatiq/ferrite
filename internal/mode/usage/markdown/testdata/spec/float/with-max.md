# Environment Variables

## Specification

### `WEIGHT`

> weighting for this node

The `WEIGHT` variable's value **MUST** be `+20.5` or less.

```bash
export WEIGHT=+20.5          # (non-normative) the maximum accepted value
export WEIGHT=-1.8715529e+38 # (non-normative)
export WEIGHT=-1.3611294e+38 # (non-normative)
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

</details>
