/*
Package dhall implements routines for deserializing and evaluating
Dhall configuration into Go data structures.

For more on the Dhall language, see https://dhall-lang.org/

If you have the following struct:

 type Message struct {
	 Name string
	 Body string
	 Time int64
 }

And a file called foo.dhall containing the following Dhall
configuration:

 { Name = "Alice", Body = "Hello", Time = 1294706395881547000 }

You can deserialize the Dhall into the struct as follows (error
handling skipped for brevity):

 var m Message
 dhallBytes, err := ioutil.ReadFile("foo.dhall")
 err = dhall.Unmarshal(dhallBytes, &m)

This version supports Dhall standard 15.0.0, except that it doesn't
support `using` directives.
*/
package dhall
