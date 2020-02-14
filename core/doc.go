/*
Package core contains the core Dhall language implementation.

If you have a Term, you can call TypeOf() to ensure it is well-typed
and to get its type.  A well-typed Term can be Eval()ed to return a
Value.  A Value represents a fully beta-normalized Dhall value.  If
you have a Value, you can use Quote() to express it as a Term.
*/
package core
