# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->

[keep a changelog]: https://keepachangelog.com/en/1.0.0/
[semantic versioning]: https://semver.org/spec/v2.0.0.html

## [Unreleased]

### Added

- Added `FileBuilder.WithMustExist()` constraint.

### Changed

- Bumped minimum Go version to v1.22.

## [1.3.0] - 2024-03-14

### Added

- Added `StringBuilder.WithMinimumLength()`, `WithMaximumLength()`, `WithLength()`.
- Added `BytesBuilder.WithMinimumLength()`, `WithMaximumLength()`, `WithLength()`.

### Fixed

- Fixed wording of error message when a variable exceeds the maximum length.

## [1.2.1] - 2023-08-14

### Changed

- Dropped support for Go v1.19, which reached end-of-life on 2023-08-08
- Removed use of the `slices` and `maps` packages from `golang.org/x/exp`. These
  package are not versioned and as such breaking changes can cause conflicts
  with other dependencies. The experimental `constraints` package is still used,
  with the expectation that these will always remain valid for use as type
  parameter constraints.

## [1.2.0] - 2023-06-12

### Added

- Added `Registry` and `NewRegistry()`
- Added `RegistryOption` and `WithDocumentationURL()`

### Changed

- Changed `WithRegistry()` to accept a `ferrite.Registry` instead of the experimental `variable.Registry` type
- Passing `WithRegistry()` to `Init()` is now additive, instead of replacing the default registry

### Removed

- Removed the experimental `maybe` and `variable` sub-packages

### Fixed

- Prevent registration of multiple environment variables with the same name

## [1.1.0] - 2023-06-06

### Added

- Added `Binary()` and `BinaryAs()` builders for binary variables
- Added `variable.LengthLimited` schema for values that have minimum or maximum lengths

### Changed

- `variable.String` now implements the new `LengthLimited` interface
- Rename the `variable.[Min|Max]LengthError.String` field to `ViolatedSchema` and change its type to `LengthLimited`

### Removed

- Removed "platform" examples from generated documentation. These are too
  opinionated and organization-specific to be included in every project's
  documentation. Instead, we will provide different export modes that can be
  used to generate Docker/Kubernetes/etc configurations, similar to the existing
  `export/dotenv` mode.

## [1.0.3] - 2023-04-20

### Changed

- `export/dotenv` mode now includes sensitive values in the output

## [1.0.2] - 2023-03-29

### Fixed

- Added `export` keyword to output of `export/dotenv` mode

### Changed

- Improve example values for `File()` builder

## [1.0.1] - 2023-03-28

### Changed

- Render a table of environment variables in Markdown documentation instead of bullet points
- Simplify language in generated documentation

## [1.0.0] - 2023-03-17

- No changes from v0.6.0.

## [0.6.0] - 2023-03-14

### Added

- Added `SeeAlso()` option
- Added `RelevantIf()` option

### Changed

- **[BC]** `SupersededBy()` now accepts a single variable set and variadic options
- **[BC]** Renamed `Input` to `VariableSet`
- Refactored internal `variable.Variable` to allow for more flexible implementations

### Removed

- **[BC]** Removed `SeeAlso()` method from all builder types, use the new `SeeAlso()` option instead

## [0.5.0] - 2023-03-12

### Added

- Added `SeeAlso()` method to all builders, which links to another variable for documentation purposes
- Added `SupersededBy()` option, which indicates which variables to use instead of a deprecated one
- Added `WithRegistry()` option to allow registration of variables in custom variable registries

### Changed

- **[BC]** Floating point variables no longer accept `NaN`, `+Inf` or `-Inf` as valid values
- **[BC]** Replaced `WithConstraintFunc()` with `WithConstraint()`
- `variable.Register()` now directly accepts a `Registry` instead of using functional options

### Removed

- Removed `variable.Option` type and `variable.WithRegistry()` option

### Fixed

- Internal errors from `strconv.ParseXXX()` functions are no longer presented to the user verbatim

## [0.4.2] - 2023-03-11

### Added

- Added `export/dotenv` mode, which renders environment variables and their
  current values in [dotenv] format

<!-- references -->

[dotenv]: https://github.com/motdotla/dotenv

### Fixed

- Fixed issue where durations were rendered with tailing zero-valued minute components

## [0.4.1] - 2023-03-10

### Fixed

- Markdown documentation is now rendered to `STDOUT` (instead of `STDERR`) as intended

## [0.4.0] - 2023-03-10

### Added

- Added `Deprecated()` method to all builders, which marks the variable as deprecated and warns the user if it is defined
- Added `StringBuilder.WithSensitiveContent()`, which hides the variable value from console output and generated documentation
- Added functional options to `Init()`

### Changed

The following changes are technically not backwards compatible from a Go
perspective, but under normal usage (as per the examples) they do not actually
alter the usage of the Ferrite API.

- **[BC]** `NetworkPort()` now returns a `NetworkPortBuilder` instead of `StringBuilder`
- **[BC]** `Required[T]` and `Optional[T]` are now structs instead of interfaces
- **[BC]** All builder types now use pointer receivers
- **[BC]** `Require()` and `Optional()` methods on all builders now use distinct option types

## [0.3.6] - 2023-03-06

### Fixed

- Fixed schema rendering of URL variables in validation mode

### Changed

- Improved algorithm for choosing which example to show in generated usage documentation

## [0.3.5] - 2023-03-06

### Fixed

- Remove use of `strings.CutPrefix()` to maintain compatibility with Go v1.19

## [0.3.4] - 2023-03-05

### Changed

- Generate Markdown documentation for the syntax of integers, floating-point values, URLs and network ports
- Other general improvements to generated Markdown documentation

### Fixed

- Trailing zeroes are no longer rendered on floating point values

## [0.3.3] - 2023-03-03

### Added

- Add `URL()`
- Generated Markdown documentation now indicates whether example values are normative or non-normative
- Generated Markdown documentation now describes the format of duration variables

## [0.3.2] - 2022-11-30

### Added

- Add `Float()`
- Add `File()`
- Add `DurationBuilder.WithMinimum()` and `WithMaximum()`
- Add `SignedBuilder.WithMinimum()` and `WithMaximum()`
- Add `UnsignedBuilder.WithMinimum()` and `WithMaximum()`

## [0.3.1] - 2022-08-24

### Added

- Add `NetworkPort()`

## [0.3.0] - 2022-08-15

### Added

- Added support for generating markdown documentation

## [0.2.0] - 2022-08-08

This is a fairly substantial refactor from the initial prototype, though usage
is largely unchanged. See the examples for details.

The two most important changes in usage are:

- `ValidateEnvironment()` has been renamed to `Init()`
- Either `Optional()` or `Required()` must be called at the end of a builder chain

## [0.1.0] - 2022-08-03

- Initial release

<!-- references -->

[unreleased]: https://github.com/dogmatiq/ferrite
[0.1.0]: https://github.com/dogmatiq/ferrite/releases/tag/v0.1.0
[0.2.0]: https://github.com/dogmatiq/ferrite/releases/tag/v0.2.0
[0.3.0]: https://github.com/dogmatiq/ferrite/releases/tag/v0.3.0
[0.3.1]: https://github.com/dogmatiq/ferrite/releases/tag/v0.3.1
[0.3.2]: https://github.com/dogmatiq/ferrite/releases/tag/v0.3.2
[0.3.3]: https://github.com/dogmatiq/ferrite/releases/tag/v0.3.3
[0.3.4]: https://github.com/dogmatiq/ferrite/releases/tag/v0.3.4
[0.3.5]: https://github.com/dogmatiq/ferrite/releases/tag/v0.3.5
[0.3.6]: https://github.com/dogmatiq/ferrite/releases/tag/v0.3.6
[0.4.0]: https://github.com/dogmatiq/ferrite/releases/tag/v0.4.0
[0.4.1]: https://github.com/dogmatiq/ferrite/releases/tag/v0.4.1
[0.4.2]: https://github.com/dogmatiq/ferrite/releases/tag/v0.4.2
[0.5.0]: https://github.com/dogmatiq/ferrite/releases/tag/v0.5.0
[0.6.0]: https://github.com/dogmatiq/ferrite/releases/tag/v0.6.0
[1.0.0]: https://github.com/dogmatiq/ferrite/releases/tag/v1.0.0
[1.0.1]: https://github.com/dogmatiq/ferrite/releases/tag/v1.0.1
[1.0.2]: https://github.com/dogmatiq/ferrite/releases/tag/v1.0.2
[1.0.3]: https://github.com/dogmatiq/ferrite/releases/tag/v1.0.3
[1.1.0]: https://github.com/dogmatiq/ferrite/releases/tag/v1.1.0
[1.2.0]: https://github.com/dogmatiq/ferrite/releases/tag/v1.2.0
[1.2.1]: https://github.com/dogmatiq/ferrite/releases/tag/v1.2.1
[1.3.0]: https://github.com/dogmatiq/ferrite/releases/tag/v1.3.0

<!-- version template
## [0.0.1] - YYYY-MM-DD

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->
