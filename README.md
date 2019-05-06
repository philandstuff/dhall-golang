# dhall-golang

Go bindings for dhall.

## Development

### Running the tests

    go test ./...

## Progress

 - [X] Type, Kind, Sort
 - [X] Variables
   - [X] de bruijn indices
   - [ ] quoted variables
 - [X] Lambdas, Pis, function application
   - [x] Alpha normalization
 - [X] Let bindings
 - [X] Type annotations
 - [X] Bools
   - [X] if
   - [x] `&&`, `||`
   - [x] `==`, `!=`
 - [X] Naturals
   - [X] `l + r` Natural addition
   - [x] `l * r` Natural multiplication
   - [ ] Natural/* standard functions
 - [X] Integers
   - [ ] Integer/toDouble and Integer/show
 - [X] Doubles
   - [ ] Double/show
 - [X] Lists
   - [x] `l # r` list append
   - [ ] List/* functions
 - [x] Text
   - [x] double quote literals
   - [x] single quote literals
   - [x] text interpolation
   - [ ] `l ++ r` text append
   - [ ] Text/show standard functions
 - [x] Optionals
   - [ ] Optional/fold and Optional/build
 - [x] Records
   - [x] `f.a`
   - [ ] `f.{ xs… }`
   - [ ] `l ∧ r`
   - [ ] `l ⫽ r`
   - [ ] `l ⩓ r`
 - [ ] Unions
   - [ ] `merge`
 - [ ] Imports
   - [x] local imports (except home-rooted paths)
   - [x] remote imports
   - [x] environment variable imports
   - [ ] import caching
   - [x] importing expressions
   - [x] importing `as Text`
   - [ ] `x ? y` alternate import operator
   - [x] `missing`
 - [X] unmarshalling into Go types
 - [ ] better errors
