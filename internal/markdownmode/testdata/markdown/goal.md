# Environment Variables

This document describes the environment variables used by `binary-name`. It is generated automatically by [dogmatiq/ferrite].

## Index

- [COMMAND_BUFFER](#COMMAND_BUFFER) — maximum number of command messages to buffer in memory
- [DEBUG](#DEBUG) — enable debug mode
- [GRPC_TIMEOUT](#GRPC_TIMEOUT) — gRPC request timeout
- [LOG_LEVEL](#LOG_LEVEL) — the minimum log level to record
- [PPROF_PORT](#PPROF_PORT) — HTTP port for serving pprof profiling data
- [READ_DSN](#READ_DSN) — database connection string for read-models
- [REDIS_SERVICE_HOST](#REDIS_SERVICE_HOST) — host for the "redis" kubernetes service
- [REDIS_SERVICE_PORT](#REDIS_SERVICE_PORT) — port for the "redis" kubernetes service

## Specification

### `COMMAND_BUFFER`

> maximum number of command messages to buffer in memory

This variable is **required**, although a default is provided.

- must be a valid [`int`]
- must be `0` or greater

```bash
export GRPC_TIMEOUT=10s   # default value
export GRPC_TIMEOUT=1m15s # example only
```

### `DEBUG`

> enable debug mode

This variable is **required**, although a default is provided.

- must be one of the values described below

```bash
export DEBUG=true
export DEBUG=false # default value
```

### `GRPC_TIMEOUT`

> gRPC request timeout

This variable is **required**, although a default is provided.

- must be a valid [`time.Duration`]
- must be `1ns` or greater

```bash
export GRPC_TIMEOUT=10s   # default value
export GRPC_TIMEOUT=1m15s # example only
```

### `LOG_LEVEL`

> the minimum log level to record

This variable is **required**, although a default is provided.

- must be one of the values described below

```bash
export LOG_LEVEL=debug # show information for developers
export LOG_LEVEL=info  # standard log messages
export LOG_LEVEL=warn  # important, but don't need individual human review
export LOG_LEVEL=error # (default) a healthy application shouldn't produce any errors
export LOG_LEVEL=fatal # the application cannot proceed
```

### `PPROF_PORT`

> HTTP port for serving pprof profiling data

This variable is **optional**,

- may be undefined
- may be an empty string
- may be defined, in which case:
  - must be a valid [`int16`]
  - must be `0` or greater

```bash
export PPROF_PORT=     # default value
export PPROF_PORT=3425 # example only
```

### `READ_DSN`

> database connection string for read-models

This variable is **required** and there is no default. **It must be set explicitly.**

```bash
export READ_DSN="???" # no meaningful example is available
```

### `REDIS_SERVICE_HOST`

> host for the "redis" kubernetes service

This variable is **required**, but it [set automatically](kubernetes) within the Kubernetes cluster.

```bash
export REDIS_SERVICE_HOST=redis.example.org # example only
```

### `REDIS_SERVICE_PORT`

> port for the "redis" kubernetes service

This variable is **required**, but it [set automatically](kubernetes) within the Kubernetes cluster.

```bash
export REDIS_SERVICE_PORT=443   # numeric port, example only
export REDIS_SERVICE_PORT=https # IANA service name, example only
```

[`time.duration`]: https://pkg.go.dev/time#ParseDuration
[`int`]: https://pkg.go.dev/builtin#int
[`int16`]: https://pkg.go.dev/builtin#int16
[kubernetes]: https://kubernetes.io/docs/concepts/services-networking/service/#environment-variables
[dogmatiq/ferrite]: https://github.com/dogmatiq/ferrite
