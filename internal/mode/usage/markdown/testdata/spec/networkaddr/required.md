# Environment Variables

## `LISTEN_ADDR`

> listen address for the HTTP server

The `LISTEN_ADDR` variable's value **MUST** be a valid network address.

```bash
export LISTEN_ADDR=192.168.0.1:8080       # (non-normative) an IPv4 address with a port
export LISTEN_ADDR='[::1]:8080'           # (non-normative) an IPv6 address with a port
export LISTEN_ADDR=host.example.org:https # (non-normative) a named host with an IANA service name
```

<details>
<summary>Network address syntax</summary>

Addresses may be specified as `<host>:<port>`, where `<host>` is a hostname or
IP address and `<port>` is a numeric port number or an IANA service name. IPv6
addresses must be enclosed in square brackets, e.g. `[::1]:8080`.

</details>
