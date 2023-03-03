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

### Automatic Usage Documentation

Ferrite can automatically generate Markdown documentation for the declared
environment variables by executing the application with the `FERRITE_MODE`
environment variable set to `usage/markdown`.

This causes the `ferrite.Init()` function to print the Markdown to `STDERR`, and
then exit the process before the application code is executed.

### Other Implementations

[Austenite](https://github.com/eloquent/austenite) is a TypeScript
library with similar features to Ferrite.
