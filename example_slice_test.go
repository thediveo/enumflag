package enumflag_test

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
)

// ① Define your new enum flag type. It can be derived from enumflag.Flag, but
// it doesn't need to be as long as it is compatible with enumflag.Flag, so
// either an int or uint.
type MooMode enumflag.Flag

// ② Define the enumeration values for FooMode.
const (
	Moo MooMode = (iota + 1) * 111
	Møø
	Mimimi
)

// ③ Map enumeration values to their textual representations (value
// identifiers).
var MooModeIds = map[MooMode][]string{
	Moo:    {"moo"},
	Møø:    {"møø"},
	Mimimi: {"mimimi"},
}

// User-defined enum flag types should be derived from "enumflag.Flag"; however
// this is not strictly necessary as long as they can be converted into the
// "enumflag.Flag" type. Actually, "enumflag.Flag" is just a fancy name for an
// "uint". In order to use such user-defined enum flags as flag slices, simply
// wrap them using enumflag.NewSlice.
func Example_slice() {
	// ④ Define your enum slice flag value.
	var moomode []MooMode
	rootCmd := &cobra.Command{
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Printf("mode is: %d=%q\n",
				moomode,
				cmd.PersistentFlags().Lookup("mode").Value.String())
		},
	}
	// ⑤ Define the CLI flag parameters for your wrapped enum slice flag.
	rootCmd.PersistentFlags().VarP(
		enumflag.NewSlice(&moomode, "mode", MooModeIds, enumflag.EnumCaseInsensitive),
		"mode", "m",
		"can be any combination of 'moo', 'møø', 'mimimi'")

	rootCmd.SetArgs([]string{"--mode", "Moo,møø"})
	_ = rootCmd.Execute()
	// Output: mode is: [111 222]="[moo,møø]"
}
