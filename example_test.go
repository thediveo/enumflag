package enumflag_test

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
)

// ① Defines a new enum flag type.
type FooMode enumflag.Flag

// ② Defines the enumeration values for our new FooMode enum flag type.
const (
	Foo FooMode = iota
	Bar
)

// ③ Implements the methods required by spf13/cobra in order to use the enum
// as a flag.
func (f *FooMode) String() string     { return enumflag.String(f) }
func (f *FooMode) Set(s string) error { return enumflag.Set(f, s) }
func (f *FooMode) Type() string       { return "foomode" }

// ④ Implements the method required by enumflag to map enum values to their
// textual identifiers.
func (f *FooMode) Enums() (interface{}, enumflag.EnumCaseSensitivity) {
	return map[FooMode][]string{
		Foo: {"foo"},
		Bar: {"bar"},
	}, enumflag.EnumCaseInsensitive
}

// New enum flag types should be derived from "enumflag.Flag"; however this is
// not strictly necessary as long as they can be converted into the
// "enumflag.Flag" type. Actually, "enumflag.Flag" is just a fancy name for an
// "uint". Enum flag types need to implement pflag's "Value" interface
// (https://godoc.org/github.com/spf13/pflag#Value), as well as the enumflag's
// "Mapper" interface. This example shows the boilerplate code, which should be
// easy to copy, paste, and adapt: simply change the Type()-returned name, and
// the Enums() returned textual enum representations and case sensitivity.
func Example() {
	var foomode FooMode // ⑤ Now use the FooMode enum flag.
	rootCmd := &cobra.Command{
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("mode is: %d=%q\n", foomode, foomode.String())
		},
	}
	// ⑥ Define the parameters for our FooMode enum flag.
	rootCmd.PersistentFlags().VarP(
		&foomode,
		"mode", "m",
		"foos the output; can be 'foo' or 'bar'")
	rootCmd.SetArgs([]string{"--mode", "bAr"})
	_ = rootCmd.Execute()
	// Output: mode is: 1="bar"
}
