# Environment Variables

## `PORT`

> listen port for the HTTP server

⚠️ The `PORT` variable is **deprecated**; its use is **NOT RECOMMENDED** as it
may be removed in a future version. If defined, the value **MUST** be a valid
network port.

```bash
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
