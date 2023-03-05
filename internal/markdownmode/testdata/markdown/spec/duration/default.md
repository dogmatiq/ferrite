# Environment Variables

## Specification

### `GRPC_TIMEOUT`

> gRPC request timeout

The `GRPC_TIMEOUT` variable **MAY** be left undefined, in which case the default
value of `10ms` is used. Otherwise, the value **MUST** be `1ns` or greater.

```bash
export GRPC_TIMEOUT=10ms # (default)
export GRPC_TIMEOUT=1ns  # (non-normative) the minimum accepted value
```

<details>
<summary>Duration syntax</summary>

Durations are specified as a sequence of decimal numbers, each with an optional
fraction and a unit suffix, such as `300ms`, `-1.5h` or `2h45m`. Supported time
units are `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`.

</details>
