# Ferrite

[![Build Status](https://github.com/dogmatiq/ferrite/actions/workflows/ci.yml/badge.svg)](https://github.com/dogmatiq/ferrite/actions/workflows/ci.yml)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/ferrite/main.svg)](https://codecov.io/github/dogmatiq/ferrite)
[![Latest Version](https://img.shields.io/github/tag/dogmatiq/ferrite.svg?label=semver)](https://semver.org)
[![Documentation](https://img.shields.io/badge/go.dev-reference-007d9c)](https://pkg.go.dev/github.com/dogmatiq/ferrite)
[![Go Report Card](https://goreportcard.com/badge/github.com/dogmatiq/ferrite)](https://goreportcard.com/report/github.com/dogmatiq/ferrite)

Ferrite is a type-safe, declarative environment variable validation system for
Go.

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

This causes the `ferrite.Init()` function to print the Markdown to `STDOUT`, and
then exit the process before the application code is executed.

### Other Implementations

[Austenite](https://github.com/eloquent/austenite) is a TypeScript
library with similar features to Ferrite.
