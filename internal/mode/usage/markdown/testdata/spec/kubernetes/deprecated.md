# Environment Variables

## Specification

### `REDIS_SERVICE_HOST`

> kubernetes "redis" service host

⚠️ The `REDIS_SERVICE_HOST` variable is **deprecated**; its use is **NOT
RECOMMENDED** as it may be removed in a future version. If defined, the value
**MUST** be a valid hostname.

It is expected that this variable will be implicitly defined by Kubernetes; it
typically does not need to be specified in the pod manifest.

```bash
export REDIS_SERVICE_HOST=foo # (non-normative)
```

#### See Also

- ~~[`REDIS_SERVICE_PORT`]~~ — ~~kubernetes "redis" service port~~ (deprecated)

### `REDIS_SERVICE_PORT`

> kubernetes "redis" service port

⚠️ The `REDIS_SERVICE_PORT` variable is **deprecated**; its use is **NOT
RECOMMENDED** as it may be removed in a future version. If defined, the value
**MUST** be a valid network port.

It is expected that this variable will be implicitly defined by Kubernetes; it
typically does not need to be specified in the pod manifest.

```bash
export REDIS_SERVICE_PORT=foo # (non-normative)
```

<details>
<summary>Network port syntax</summary>

Ports may be specified as a numeric value no greater than `65535`.
Alternatively, a service name can be used. Service names are resolved against
the system's service database, typically located in the `/etc/service` file on
UNIX-like systems. Standard service names are published by IANA.

</details>

#### See Also

- ~~[`REDIS_SERVICE_HOST`]~~ — ~~kubernetes "redis" service host~~ (deprecated)

<!-- references -->

[`redis_service_host`]: #REDIS_SERVICE_HOST
[`redis_service_port`]: #REDIS_SERVICE_PORT
