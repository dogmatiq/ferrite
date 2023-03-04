# Environment Variables

## Specification

### `PPROF_PORT`

> HTTP port for serving pprof profiling data

This variable **MAY** be set to a `uint16` value.
If left undefined the default value of `8080` is used.

```bash
export PPROF_PORT=8080 # (default)
```
