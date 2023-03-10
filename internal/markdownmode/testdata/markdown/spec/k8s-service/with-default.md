# Environment Variables

## Specification

### `REDIS_SERVICE_HOST`

> kubernetes "redis" service host

The `REDIS_SERVICE_HOST` variable **MAY** be left undefined, in which case the
default value of `redis.example.org` is used. Otherwise, the value **MUST** be a
valid hostname.

It is expected that this variable will be implicitly defined by Kubernetes; it
typically does not need to be specified in the pod manifest.

```bash
export REDIS_SERVICE_HOST=redis.example.org # (default)
```

### `REDIS_SERVICE_PORT`

> kubernetes "redis" service port

The `REDIS_SERVICE_PORT` variable **MAY** be left undefined, in which case the
default value of `6379` is used. Otherwise, the value **MUST** be a valid
network port.

It is expected that this variable will be implicitly defined by Kubernetes; it
typically does not need to be specified in the pod manifest.

```bash
export REDIS_SERVICE_PORT=6379 # (default)
```

<details>
<summary>Network port syntax</summary>

Ports may be specified as a numeric value no greater than `65535`.
Alternatively, a service name can be used. Service names are resolved against
the system's service database, typically located in the `/etc/service` file on
UNIX-like systems. Standard service names are published by IANA.

</details>
