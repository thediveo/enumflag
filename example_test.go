package enumflag_test

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"
)

// ① Define your new enum flag type. It can be derived from enumflag.Flag,
// but it doesn't need to be as long as it satisfies constraints.Integer.
type FooMode enumflag.Flag

// ② Define the enumeration values for FooMode.
const (
	Foo FooMode = iota
	Bar
)

// ③ Map enumeration values to their textual representations (value
// identifiers).
var FooModeIds = map[FooMode][]string{
	Foo: {"foo"},
	Bar: {"bar"},
}

func Example() {
	// ④ Define your enum flag value.
	var foomode FooMode
	rootCmd := &cobra.Command{
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Printf("mode is: %d=%q\n",
				foomode,
				cmd.PersistentFlags().Lookup("mode").Value.String())
		},
	}
	// ⑤ Define the CLI flag parameters for your wrapped enum flag.
	rootCmd.PersistentFlags().VarP(
		enumflag.New(&foomode, "mode", FooModeIds, enumflag.EnumCaseInsensitive),
		"mode", "m",
		"foos the output; can be 'foo' or 'bar'")

	rootCmd.SetArgs([]string{"--mode", "bAr"})
	_ = rootCmd.Execute()
	// Output: mode is: 1="bar"
}
