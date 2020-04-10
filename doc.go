/*

Package enumflag supplements the Golang CLI flag handling packages spf13/cobra
and spf13/pflag with enumeration flags.

Enumeration flags allow users to specify flag values based on a set of
specific enumeration text values, such as "--mode=foo" or "--mode=bar".
Internally, the application programmer then simply deals with enumeration
values in form of uints (or ints), without needing to parse and validate
enumeration flags herself.

*/
package enumflag
