# Changelog

## [Unreleased]
[Unreleased]: https://github.com/philandstuff/dhall-golang/compare/v4.0.0...HEAD

## [4.0.0] - 2020-06-16
[4.0.0]: https://github.com/philandstuff/dhall-golang/compare/v3.0.0...v4.0.0

This brings dhall-golang up to version 17.0.0 of the Dhall standard.
Again the standard had breaking changes, so this release is a major
version bump.

Thanks to @lisael for their contributions to this release.

### Breaking changes

 * Language changes:
   * [Remove Optional/build and Optional/fold](https://github.com/dhall-lang/dhall-lang/pull/1014)

### Added

 * Language changes:
   * [Allow quoted labels to be empty](https://github.com/dhall-lang/dhall-lang/pull/980)

### Fixed

 * Fix potential stack overflow in typechecker (#40)
    * When typechecking certain pathological expressions, the
      typechecker would get into an infinite loop until it exhausted
      the stack.
 * Fix error messages when `x === y` fails to typecheck (#39)

## [3.0.0] - 2020-05-11
[3.0.0]: https://github.com/philandstuff/dhall-golang/compare/v2.0.0...v3.0.0

This brings dhall-golang up to version 16.0.0 of the Dhall standard.
As the standard had breaking changes, this release is a major version
bump.

### Breaking changes

 * Language changes:
     * [Adjust precedence of `===` and `with`](https://github.com/dhall-lang/dhall-lang/pull/954)
     * [Update encoding of floating point values to RFC7049bis](https://github.com/dhall-lang/dhall-lang/pull/958)

### New features

 * Language features:
     * [Allow unions with mixed kinds](https://github.com/dhall-lang/dhall-lang/pull/957)

### Bug fixes

 * We now save fully-alpha-normalized expressions to the cache (#31)
 * We now check the hash of expressions fetched from the cache (#32)

## [2.0.0] - 2020-04-17
[2.0.0]: https://github.com/philandstuff/dhall-golang/compare/v1.0.0...v2.0.0

This brings dhall-golang up to version 15.0.0 of the Dhall standard.
As the standard had breaking changes, this release is a major version
bump.

### Breaking changes

 - added `with` keyword (technically breaking since you can no longer
   use `with` as an identifier)

### Added

 - added record puns (ie `{ x }` is now shorthand for `{ x = x }`)
 - added UnmarshalFile function (#25)

### Changed

 - Unmarshal() and Decode() will check a Dhall function matches the
   given Go type before decoding (#23)
 - imports are now evaluated at import time (#27)

### Fixed

 - fixed bug in evaluation of `merge` (3171f34)

## [1.0.0] - 2020-03-15
[1.0.0]: https://github.com/philandstuff/dhall-golang/compare/v1.0.0-rc.4...v1.0.0

No changes from 1.0.0-rc.4.

## [1.0.0-rc.4] - 2020-03-06
[1.0.0-rc.4]: https://github.com/philandstuff/dhall-golang/compare/v1.0.0-rc.3...v1.0.0-rc.4

### Breaking changes

 - dhall.Decode() now returns an error instead of panicking (#18)

### Changed

 - regenerate parser from mna/pigeon master (#16)
 - support for unmarshalling into pointer types (#17)
 - better encoding of Optional types (#19)

## [1.0.0-rc.3] - 2020-02-23

[1.0.0-rc.3]: https://github.com/philandstuff/dhall-golang/compare/v1.0.0-rc.2...v1.0.0-rc.3

Another release candidate.  A few more breaking changes, though less
drastic than rc.2 was.  Things are slowly stabilising.

Thanks to @Duncaen for his contribution to this release.

### Breaking changes

 - core.Pi.Range has been renamed to core.Pi.Codomain (#12)
 - core.TextLit has been removed and replaced with core.PlainTextLit.
   There is no longer a (public) interpolated text Value type.
 - struct tags are now `dhall` not `json` (#15)

### Fixed

 - fixed a parser bug related to single quote strings (#1)

### Changed

 - faster parser by approx 30% (#11)

## [1.0.0-rc.2] - 2020-02-16

[1.0.0-rc.2]: https://github.com/philandstuff/dhall-golang/compare/v1.0.0-rc.1...v1.0.0-rc.2

Another release candidate.  As promised, the `core` package is still
in flux and has undergone a huge refactor in this release.  Along with
that, the godoc has been vastly improved, and a new README has been
written to replace the previous scrappy development notes.

Also, this brings dhall-golang up to version 14.0.0 of the language
standard.

### Breaking changes

 - refactoring of the `core` package
   - moved `core.Term` and implementations to new package `term`
   - removed `core.AppValue` from public interface
   - renamed various types to remove `-Val` and `-Term` suffixes
 - (from Dhall 14.0.0): decimal Natural literals can no longer have
   leading 0 digits

### Changed

 - dhall-golang now supports [version 14.0.0][dhall-14.0.0] of the
   language.

[dhall-14.0.0]: https://github.com/dhall-lang/dhall-lang/releases/tag/v14.0.0

### Fixed

 - `dhall.Unmarshal()` now resolves imports and typechecks before
   evaluating

### Added

  - added `core.ListOf`, `core.OptionalOf`, `core.NoneOf` Value types
    to represent `List a`, `Optional a` and `None a` Values
    respectively

## [1.0.0-rc.1] - 2020-02-09

### Changed

 - Fixed a compile error in cbor.go :/

[1.0.0-rc.1]: https://github.com/philandstuff/dhall-golang/compare/v1.0.0-rc.0...v1.0.0-rc.1

## [1.0.0-rc.0] - 2020-02-09

First release candidate.  Note that some things are still in flux and
subject to change:

 - The `dhall` package is stable and will not have any breaking
   changes.  In particular, `dhall.Decode` and `dhall.Unmarshal` will
   not have any breaking changes before a v1.0.0 release.
 - The `parser` package is also stable and will not have any breaking
   changes.
 - The `core` package is still subject to change: in particular, names
   which are currently exported may be unexported before a v1.0.0
   release.

### Added

- Core Dhall functionality:
  - Parse Dhall source to Terms
  - Resolve Dhall imports
  - Use Dhall cache for imports
  - Typecheck Dhall Terms
  - Evaluate Dhall Terms to Values
  - Marshalling/unmarshalling to CBOR format
- Go bindings:
  - dhall.Decode to decode a Dhall Value into a Go variable
  - dhall.Unmarshal as a convenience all-in-one
    Dhall-source-to-Go-variable function

[1.0.0-rc.0]: https://github.com/philandstuff/dhall-golang/releases/tag/v1.0.0-rc.0
