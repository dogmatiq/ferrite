# Environment Variables

| Name        | Usage    | Description                          |
| ----------- | -------- | ------------------------------------ |
| [`DEBUG`]   | optional | enable or disable debugging features |
| [`VERBOSE`] | optional | enable verbose logging               |

## `DEBUG`

> enable or disable debugging features

The `DEBUG` variable **MAY** be left undefined. Otherwise, the value **MUST** be
either `true` or `false`.

```bash
export DEBUG=true
export DEBUG=false
```

### See Also

- [`VERBOSE`] — enable verbose logging

## `VERBOSE`

> enable verbose logging

The `VERBOSE` variable **MAY** be left undefined. Otherwise, the value **MUST**
be either `true` or `false`.

```bash
export VERBOSE=true
export VERBOSE=false
```

<!-- references -->

[`debug`]: #DEBUG
[`verbose`]: #VERBOSE
