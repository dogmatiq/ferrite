# Environment Variables

This document describes the environment variables used by `<app>`.

| Name      | Optionality | Description                          |
| --------- | ----------- | ------------------------------------ |
| [`DEBUG`] | optional    | enable or disable debugging features |

⚠️ `<app>` may consume other undocumented environment variables. This document
only shows variables declared using [Ferrite].

## Specification

All environment variables described below must meet the stated requirements.
Otherwise, `<app>` prints usage information to `STDERR` then exits.
**Undefined** variables and **empty** values are equivalent.

The key words **MUST**, **MUST NOT**, **REQUIRED**, **SHALL**, **SHALL NOT**,
**SHOULD**, **SHOULD NOT**, **RECOMMENDED**, **MAY**, and **OPTIONAL** in this
document are to be interpreted as described in [RFC 2119].

### `DEBUG`

> enable or disable debugging features

The `DEBUG` variable **MAY** be left undefined. Otherwise, the value **MUST** be
either `true` or `false`.

```bash
export DEBUG=true
export DEBUG=false
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
            - name: DEBUG # enable or disable debugging features (optional)
              value: "false"
```

Alternatively, the environment variables can be defined within a [config map][kubernetes config map]
then referenced from a deployment manifest using `configMapRef`.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: example-config-map
data:
  DEBUG: "false" # enable or disable debugging features (optional)
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
      DEBUG: "false" # enable or disable debugging features (optional)
```

</details>

<!-- references -->

[`debug`]: #DEBUG
[docker service]: https://docs.docker.com/compose/environment-variables/#set-environment-variables-in-containers
[ferrite]: https://github.com/dogmatiq/ferrite
[kubernetes config map]: https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#configure-all-key-value-pairs-in-a-configmap-as-container-environment-variables
[kubernetes container]: https://kubernetes.io/docs/tasks/inject-data-application/define-environment-variable-container/#define-an-environment-variable-for-a-container
[rfc 2119]: https://www.rfc-editor.org/rfc/rfc2119.html
