# Environment Variables

## Specification

### `WEIGHT`

> weighting for this node

The `WEIGHT` variable's value **MUST** be `10` or greater.

```bash
export WEIGHT=10    # (non-normative) the minimum accepted value
export WEIGHT=29496 # (non-normative)
export WEIGHT=39325 # (non-normative)
```

<details>
<summary>Unsigned integer syntax</summary>

Unsigned integers can only be specified using decimal (base-10) notation. A
leading sign (`+` or `-`) is not supported and **MUST NOT** be specified.

Internally, the `WEIGHT` variable is represented using an unsigned 16-bit
integer type (`uint16`); any value that overflows this data-type is invalid.

</details>
