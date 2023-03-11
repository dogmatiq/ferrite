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

- [`PORT`] — an environment variable that has a default value

## Specification

### `PORT`

> an environment variable that has a default value

The `PORT` variable **MAY** be left undefined, in which case the default value
of `ftp` is used. Otherwise, the value **MUST** be a valid network port.

```bash
export PORT=ftp   # (default)
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
            - name: PORT # an environment variable that has a default value
              value: ftp
```

Alternatively, the environment variables can be defined within a [config map][kubernetes config map]
then referenced a deployment manifest using `configMapRef`.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: example-config-map
data:
  PORT: ftp # an environment variable that has a default value
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
      PORT: ftp # an environment variable that has a default value
```

</details>

<!-- references -->

[docker service]: https://docs.docker.com/compose/environment-variables/#set-environment-variables-in-containers
[ferrite]: https://github.com/dogmatiq/ferrite
[kubernetes config map]: https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#configure-all-key-value-pairs-in-a-configmap-as-container-environment-variables
[kubernetes container]: https://kubernetes.io/docs/tasks/inject-data-application/define-environment-variable-container/#define-an-environment-variable-for-a-container
[`port`]: #PORT
[rfc 2119]: https://www.rfc-editor.org/rfc/rfc2119.html
