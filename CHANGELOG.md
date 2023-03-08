# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->

[keep a changelog]: https://keepachangelog.com/en/1.0.0/
[semantic versioning]: https://semver.org/spec/v2.0.0.html

## [Unreleased]

### Added

- Added `InitOption` type
- Added `VariableOption` type
- Added `Sensitive()` option, which indicates that a variable's contents may contain sensitive information

### Changed

The following changes are technically not backwards compatible from a Go
perspective, but under normal usage (as per the examples) they do not
actually alter the usage of the Ferrite API.

- **[BC]** Change `Init()` to accept options
- **[BC]** Change `NetworkPort()` to use the new `NetworkPortBuilder` instead of `StringBuilder`
- **[BC]** Changed `Required[T]` and `Optional[T]` from structs to interfaces
- **[BC]** Changed all "builder" types to use pointer receivers
- **[BC]** `Require()` and `Optional()` methods on builders no longer accept `variable.RegisterOption` values

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

<!-- version template
## [0.0.1] - YYYY-MM-DD

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->
