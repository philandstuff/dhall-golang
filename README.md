# dhall-golang

Go bindings for dhall.

## Development

### Running the tests

    go test ./...

## Progress

 - [X] Type, Kind, Sort
 - [X] Variables
   - [X] de bruijn indices
 - [X] Lambdas, Pis, function application
   - [ ] Alpha normalization
 - [X] Let bindings
 - [X] Type annotations
 - [X] Bools
   - [X] if
   - [ ] `&&`, `||`
   - [ ] `==`, `!=`
 - [X] Naturals
   - [X] `l + r` Natural addition
   - [x] `l * r` Natural multiplication
   - [ ] Natural/* standard functions
 - [X] Integers
   - [ ] Integer/toDouble and Integer/show
 - [X] Doubles
   - [ ] Double/show
 - [X] Lists
   - [ ] `l # r` list append
   - [ ] List/* functions
 - [ ] Text
   - [ ] text interpolation
   - [ ] multiline strings
   - [ ] `l ++ r` text append
   - [ ] Text/show standard functions
 - [ ] Optionals
   - [ ] Optional/fold and Optional/build
 - [x] Records
   - [ ] `f.a`
   - [ ] `f.{ xs… }`
   - [ ] `l ∧ r`
   - [ ] `l ⫽ r`
   - [ ] `l ⩓ r`
 - [ ] Unions
   - [ ] `merge`
 - [ ] Imports
   - [ ] local imports
   - [ ] remote imports
   - [ ] environment variable imports
 - [X] unmarshalling into Go types
 - [ ] better errors
