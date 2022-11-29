# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->

[keep a changelog]: https://keepachangelog.com/en/1.0.0/
[semantic versioning]: https://semver.org/spec/v2.0.0.html

## [Unreleased]

### Added

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

<!-- version template
## [0.0.1] - YYYY-MM-DD

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->
