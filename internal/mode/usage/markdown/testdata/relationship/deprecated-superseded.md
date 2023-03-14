# Environment Variables

## Index

- ~~[`BIND_ADDRESS`]~~ — ~~listen address for the HTTP server~~ (deprecated)
- [`BIND_HOST`] — listen host for the HTTP server
- [`BIND_PORT`] — listen port for the HTTP server
- [`BIND_VERSION`] — IP version for the HTTP server

## Specification

### `BIND_ADDRESS`

> listen address for the HTTP server

⚠️ The `BIND_ADDRESS` variable is **deprecated**; its use is **NOT RECOMMENDED**
as it may be removed in a future version. [`BIND_HOST`], [`BIND_PORT`] and
[`BIND_VERSION`] **SHOULD** be used instead.

```bash
export BIND_ADDRESS=0.0.0.0:8080 # (default)
```

### `BIND_HOST`

> listen host for the HTTP server

The `BIND_HOST` variable **MAY** be left undefined, in which case the default
value of `0.0.0.0` is used.

```bash
export BIND_HOST=0.0.0.0 # (default)
```

### `BIND_PORT`

> listen port for the HTTP server

The `BIND_PORT` variable **MAY** be left undefined, in which case the default
value of `8080` is used. Otherwise, the value **MUST** be a valid network port.

```bash
export BIND_PORT=8080  # (default)
export BIND_PORT=8000  # (non-normative) a port commonly used for private web servers
export BIND_PORT=https # (non-normative) the IANA service name that maps to port 443
```

<details>
<summary>Network port syntax</summary>

Ports may be specified as a numeric value no greater than `65535`.
Alternatively, a service name can be used. Service names are resolved against
the system's service database, typically located in the `/etc/service` file on
UNIX-like systems. Standard service names are published by IANA.

</details>

### `BIND_VERSION`

> IP version for the HTTP server

The `BIND_VERSION` variable **MAY** be left undefined, in which case the default
value of `4` is used.

```bash
export BIND_VERSION=4 # (default)
```

<!-- references -->

[`bind_address`]: #BIND_ADDRESS
[`bind_host`]: #BIND_HOST
[`bind_port`]: #BIND_PORT
[`bind_version`]: #BIND_VERSION
