package enumflag_test

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thediveo/enumflag/v2"
)

// ① Define your new enum flag type. It can be derived from enumflag.Flag,
// but it doesn't need to be as long as it satisfies comparable.
type BarMode enumflag.Flag

// ② Define the enumeration values for BarMode.
const (
	NoDefault         = iota // optional definition for "no default" zero value
	Barr      BarMode = iota
	Barz
)

// ③ Map enumeration values to their textual representations (value
// identifiers).
var BarModeIds = map[BarMode][]string{
	// ...do NOT include/map the "no default" zero value!
	Barr: {"barr"},
	Barz: {"barz"},
}

func Example_no_default_value() {
	// ④ Define your enum flag value.
	var barmode BarMode
	rootCmd := &cobra.Command{
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Printf("mode is: %d=%q\n",
				barmode,
				cmd.PersistentFlags().Lookup("mode").Value.String())
		},
	}
	// ⑤ Define the CLI flag parameters for your wrapped enum flag.
	rootCmd.PersistentFlags().VarP(
		enumflag.NewWithoutDefault(&barmode, "mode", BarModeIds, enumflag.EnumCaseInsensitive),
		"mode", "m",
		"bars the output; can be 'barr' or 'barz'")

	// now cobra's help won't render the default enum value identifier anymore...
	_ = rootCmd.Help()

	_ = rootCmd.Execute()

	// Output:
	// Usage:
	//    [flags]
	//
	// Flags:
	//   -m, --mode mode   bars the output; can be 'barr' or 'barz'
	// mode is: 0=""
}
