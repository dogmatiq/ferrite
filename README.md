<div align="center">

# Ferrite

A type-safe, declarative environment variable validation system for Go.

[![Documentation](https://img.shields.io/badge/go.dev-documentation-007d9c?&style=for-the-badge)](https://pkg.go.dev/github.com/dogmatiq/ferrite)
[![Latest Version](https://img.shields.io/github/tag/dogmatiq/ferrite.svg?&style=for-the-badge&label=semver)](https://github.com/dogmatiq/ferrite/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/dogmatiq/ferrite/ci.yml?style=for-the-badge&branch=main)](https://github.com/dogmatiq/ferrite/actions/workflows/ci.yml)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/ferrite/main.svg?style=for-the-badge)](https://codecov.io/github/dogmatiq/ferrite)

</div>

## Getting Started

> This example demonstrates how to declare an environment variable that
> produces a `time.Duration` value, but Ferrite supports many different variable
> types, as described in the
> [examples](https://pkg.go.dev/github.com/dogmatiq/ferrite#pkg-examples).

First, describe the application's environment variables using Ferrite's
"builder" interfaces. This is typically done at the package-scope of the `main`
package.

```go
var httpTimeout = ferrite.
    Duration("HTTP_TIMEOUT", "the maximum duration of each HTTP request").
    WithDefault(10 * time.Second).
    Required()
```

Next, initialize Ferrite in the application's `main()` function before any other
code is executed so that it may halt execution when the environment is invalid.

```go
func main() {
    ferrite.Init()

    // existing application logic ...
}
```

Finally, read the environment variable's value by calling the `Value()` method.

```go
timeout := httpTimeout.Value()
ctx, cancel := context.WithTimeout(context.Background(), timeout)
defer cancel()

// do HTTP request ...
```

## Modes of Operation

By default, calling `Init()` operates in "validation" mode. There are several
other modes that can be used to gain insight into the application's use of
environment variables.

Modes are selected by setting the `FERRITE_MODE` environment variable.

### `validate` mode

This is the default mode. If one or more environment variables are invalid, this
mode renders a description of all declared environment variables and their
associated values and validation failures to `STDERR`, then exits the process
with a non-zero exit code.

It also shows warnings if deprecated environment variables are used.

### `usage/markdown` mode

This mode renders Markdown documentation about the environment variables to
`STDOUT`. The output is designed to be included in the application's `README.md`
file or a similar file.

### `export/dotenv` mode

This mode renders environment variables to `STDOUT` in a format suitable for use
with tools like [`dotenv`](https://github.com/motdotla/dotenv) and the
[`env_file`](https://docs.docker.com/compose/compose-file/#env_file) directive
in Docker compose files.

## Other Implementations

[Austenite](https://github.com/eloquent/austenite) is a TypeScript
library with similar features to Ferrite.
