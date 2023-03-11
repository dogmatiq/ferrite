# Environment Variables

This document describes the environment variables used by `<app>`.

If any of the environment variable values do not meet the requirements herein,
the application will print usage information to `STDERR` then exit with a
non-zero exit code. Please note that **undefined** variables and **empty**
values are considered equivalent.

⚠️ The application may consume other undocumented environment variables; this
document only shows those variables declared using [Ferrite].

The key words **MUST**, **MUST NOT**, **REQUIRED**, **SHALL**, **SHALL NOT**,
**SHOULD**, **SHOULD NOT**, **RECOMMENDED**, **MAY**, and **OPTIONAL** in this
document are to be interpreted as described in [RFC 2119].

## Index

- [`DEBUG`](#DEBUG) — enable or disable debugging features
- [`VERBOSE`](#VERBOSE) — enable verbose logging

## Specification

### `DEBUG`

> enable or disable debugging features

The `DEBUG` variable **MAY** be left undefined. Otherwise, the value **MUST** be
either `true` or `false`.

```bash
export DEBUG=true
export DEBUG=false
```

#### See Also

- [`VERBOSE`](#VERBOSE) — enable verbose logging

### `VERBOSE`

> enable verbose logging

The `VERBOSE` variable **MAY** be left undefined. Otherwise, the value **MUST**
be either `true` or `false`.

```bash
export VERBOSE=true
export VERBOSE=false
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
            - name: DEBUG # enable or disable debugging features
              value: "false"
            - name: VERBOSE # enable verbose logging
              value: "false"
```

Alternatively, the environment variables can be defined within a [config map][kubernetes config map]
then referenced a deployment manifest using `configMapRef`.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: example-config-map
data:
  DEBUG: "false" # enable or disable debugging features
  VERBOSE: "false" # enable verbose logging
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
      DEBUG: "false" # enable or disable debugging features
      VERBOSE: "false" # enable verbose logging
```

</details>

<!-- references -->

[docker service]: https://docs.docker.com/compose/environment-variables/#set-environment-variables-in-containers
[ferrite]: https://github.com/dogmatiq/ferrite
[kubernetes config map]: https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#configure-all-key-value-pairs-in-a-configmap-as-container-environment-variables
[kubernetes container]: https://kubernetes.io/docs/tasks/inject-data-application/define-environment-variable-container/#define-an-environment-variable-for-a-container
[rfc 2119]: https://www.rfc-editor.org/rfc/rfc2119.html
