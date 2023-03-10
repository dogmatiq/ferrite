# Environment Variables

## Specification

### `GRPC_TIMEOUT`

> gRPC request timeout

The `GRPC_TIMEOUT` variable's value **MUST** be `1ns` or greater.

```bash
export GRPC_TIMEOUT=1ns                      # (non-normative) the minimum accepted value
export GRPC_TIMEOUT=1152921h30m16.584649216s # (non-normative)
export GRPC_TIMEOUT=1537228h40m22.11286528s  # (non-normative)
```

<details>
<summary>Duration syntax</summary>

Durations are specified as a sequence of decimal numbers, each with an optional
fraction and a unit suffix, such as `300ms`, `-1.5h` or `2h45m`. Supported time
units are `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`.

</details>
