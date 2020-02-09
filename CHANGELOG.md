# Changelog

## [Unreleased]

[Unreleased]: https://github.com/philandstuff/dhall-golang/compare/v1.0.0-rc.0...HEAD

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
