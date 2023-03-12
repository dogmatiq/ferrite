# Environment Variables

This document describes the environment variables used by `<app>`.

If any of the environment variable values do not meet the requirements herein,
the application will print usage information to `STDERR` then exit with a
non-zero exit code. Please note that **undefined** variables and **empty**
values are considered equivalent.

⚠️ This document includes **non-normative** example values. While these values
are syntactically correct, they may not be meaningful to this application.

⚠️ The application may consume other undocumented environment variables; this
document only shows those variables declared using [Ferrite].

The key words **MUST**, **MUST NOT**, **REQUIRED**, **SHALL**, **SHALL NOT**,
**SHOULD**, **SHOULD NOT**, **RECOMMENDED**, **MAY**, and **OPTIONAL** in this
document are to be interpreted as described in [RFC 2119].

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

## Usage Examples

<details>
<summary>Kubernetes</summary>

This example shows how to define the environment variables needed by `<app>`
on a [Kubernetes container] within a Kubenetes deployment manifest.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: example-deployment
spec:
  template:
    spec:
      containers:
        - name: example-container
          env:
            - name: BIND_ADDRESS # listen address for the HTTP server (deprecated)
              value: 0.0.0.0:8080
            - name: BIND_HOST # listen host for the HTTP server (defaults to 0.0.0.0)
              value: 0.0.0.0
            - name: BIND_PORT # listen port for the HTTP server (defaults to 8080)
              value: "8080"
            - name: BIND_VERSION # IP version for the HTTP server (defaults to 4)
              value: "4"
```

Alternatively, the environment variables can be defined within a [config map][kubernetes config map]
then referenced a deployment manifest using `configMapRef`.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: example-config-map
data:
  BIND_ADDRESS: 0.0.0.0:8080 # listen address for the HTTP server (deprecated)
  BIND_HOST: 0.0.0.0 # listen host for the HTTP server (defaults to 0.0.0.0)
  BIND_PORT: "8080" # listen port for the HTTP server (defaults to 8080)
  BIND_VERSION: "4" # IP version for the HTTP server (defaults to 4)
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: example-deployment
spec:
  template:
    spec:
      containers:
        - name: example-container
          envFrom:
            - configMapRef:
                name: example-config-map
```

</details>

<details>
<summary>Docker</summary>

This example shows how to define the environment variables needed by `<app>`
when running as a [Docker service] defined in a Docker compose file.

```yaml
service:
  example-service:
    environment:
      BIND_ADDRESS: 0.0.0.0:8080 # listen address for the HTTP server (deprecated)
      BIND_HOST: 0.0.0.0 # listen host for the HTTP server (defaults to 0.0.0.0)
      BIND_PORT: "8080" # listen port for the HTTP server (defaults to 8080)
      BIND_VERSION: "4" # IP version for the HTTP server (defaults to 4)
```

</details>

<!-- references -->

[`bind_address`]: #BIND_ADDRESS
[`bind_host`]: #BIND_HOST
[`bind_port`]: #BIND_PORT
[`bind_version`]: #BIND_VERSION
[docker service]: https://docs.docker.com/compose/environment-variables/#set-environment-variables-in-containers
[ferrite]: https://github.com/dogmatiq/ferrite
[kubernetes config map]: https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#configure-all-key-value-pairs-in-a-configmap-as-container-environment-variables
[kubernetes container]: https://kubernetes.io/docs/tasks/inject-data-application/define-environment-variable-container/#define-an-environment-variable-for-a-container
[rfc 2119]: https://www.rfc-editor.org/rfc/rfc2119.html
