# Environment Variables

## Specification

### `WEIGHT`

> weighting for this node

The `WEIGHT` variable **SHOULD** be left undefined, in which case the default
value of `900` is used. Otherwise, the value **MUST** be a non-negative whole
number.

⚠️ This variable is **deprecated**; its use is discouraged as it may be removed
in a future version.

```bash
export WEIGHT=900 # (default)
```

<details>
<summary>Unsigned integer syntax</summary>

Unsigned integers can only be specified using decimal (base-10) notation. A
leading sign (`+` or `-`) is not supported and **MUST NOT** be specified.

Internally, the `WEIGHT` variable is represented using an unsigned 16-bit
integer type (`uint16`); any value that overflows this data-type is invalid.

</details>
