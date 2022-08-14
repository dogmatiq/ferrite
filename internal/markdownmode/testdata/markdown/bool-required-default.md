# Environment Variables

This document describes the environment variables used by `<app>`.

The application may consume other undocumented environment variables; this
document only shows those variables defined using [dogmatiq/ferrite].

## Index

- [`DEBUG`](#DEBUG) — enable or disable debugging features

## Specification

### `DEBUG`

> enable or disable debugging features

This variable is **required**, although a default is provided.
It must be one of the values shown in the examples below.

```bash
export DEBUG=true
export DEBUG=false # default value
```

<details>
<summary>Usage Examples</summary>

#### Kubernetes Container

```yaml
env:
  - name: DEBUG
    value: "true"
```

#### Kubernetes Config Map

```yaml
data:
  DEBUG: "true"
```

#### Docker Compose / Stack

```yaml
environment:
  DEBUG: "true"
```

#### GitHub Actions Workflow

```yaml
env:
  DEBUG: "true"
```

</details>

<!-- references -->

[dogmatiq/ferrite]: https://github.com/dogmatiq/ferrite
