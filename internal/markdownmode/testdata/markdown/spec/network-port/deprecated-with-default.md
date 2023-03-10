# Environment Variables

## Specification

### `PORT`

> listen port for the HTTP server

The `PORT` variable **SHOULD** be left undefined, in which case the default
value of `8080` is used. Otherwise, the value **MUST** be a valid network port.

⚠️ This variable is **deprecated**; its use is discouraged as it may be removed
in a future version.

```bash
export PORT=8080  # (default)
export PORT=8000  # (non-normative) a port commonly used for private web servers
export PORT=https # (non-normative) the IANA service name that maps to port 443
```

<details>
<summary>Network port syntax</summary>

Ports may be specified as a numeric value no greater than `65535`.
Alternatively, a service name can be used. Service names are resolved against
the system's service database, typically located in the `/etc/service` file on
UNIX-like systems. Standard service names are published by IANA.

</details>
