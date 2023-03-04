# Environment Variables

## Specification

### `GRPC_TIMEOUT`

> gRPC request timeout

This variable **MUST** be set to `1ns` or greater.

Valid durations are a sequence of decimal numbers, each with an optional
fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time
units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

```bash
export GRPC_TIMEOUT=640511h56m49.213693952s # (non-normative)
```
