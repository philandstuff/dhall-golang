# dhall-golang

Go bindings for dhall.

## Development

### Running the tests

    go test ./...

## Progress

 - [X] Type, Kind, Sort
 - [X] Variables
   - [X] de bruijn indices
   - [x] quoted variables
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
   - [x] Natural/* standard functions
 - [X] Integers
   - [x] Integer/toDouble and Integer/show
 - [X] Doubles
   - [x] Double/show (it exists but doesn't pass tests)
 - [X] Lists
   - [x] `l # r` list append
   - [x] List/* functions
 - [x] Text
   - [x] double quote literals
   - [x] single quote literals
   - [x] text interpolation
   - [x] `l ++ r` text append
   - [x] Text/show standard functions
 - [x] Optionals
   - [x] Optional/fold and Optional/build
 - [x] Records
   - [x] `f.a`
   - [x] `f.{ xs… }`
   - [x] `f.(s)`
   - [x] `l ∧ r`
   - [x] `l ⫽ r`
   - [x] `l ⩓ r`
 - [x] Unions
   - [x] types
   - [x] constructors
   - [x] `merge`
 - [ ] Imports
   - [x] local imports (except home-rooted paths)
   - [x] remote imports
   - [x] environment variable imports
   - [ ] `using ./headers`
   - [x] import caching
   - [x] importing expressions
   - [x] importing `as Text`
   - [x] `x ? y` alternate import operator
   - [x] `missing`
 - [X] unmarshalling into Go types
 - [ ] better errors
 - [ ] better godoc
