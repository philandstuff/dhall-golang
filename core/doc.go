/*
Package core contains the core Dhall language implementation.

The main types are Term and Value.  A Term represents an arbitrary
Dhall expression, while a Value represents a fully beta-normalized
Dhall value.  Eval() takes Terms and normalizes them to Values;
Quote() takes Values and expresses them as Terms.
*/
package core
