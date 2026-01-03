package enumflag_test

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"
)

func Example_external() {
	// ①+② skip "define your own enum flag type" and enumeration values, as we
	// already have a 3rd party one.

	// ③ Map 3rd party enumeration values to their textual representations
	var LoglevelIds = map[slog.Level][]string{
		slog.LevelDebug: {"debug"},
		slog.LevelInfo:  {"info"},
		slog.LevelWarn:  {"warning", "warn"},
		slog.LevelError: {"error"},
	}

	// ④ Define your enum flag value and set the your logging default value.
	var loglevel = slog.LevelWarn

	rootCmd := &cobra.Command{
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Printf("logging level is: %d=%q\n",
				loglevel,
				cmd.PersistentFlags().Lookup("log").Value.String())
		},
	}

	// ⑤ Define the CLI flag parameters for your wrapped enum flag.
	rootCmd.PersistentFlags().Var(
		enumflag.New(&loglevel, "log", LoglevelIds, enumflag.EnumCaseInsensitive),
		"log",
		"sets logging level; can be 'trace', 'debug', 'info', 'warn', 'error', 'fatal', 'panic'")

	_ = rootCmd.Execute()
	rootCmd.SetArgs([]string{"--log", "debug"})
	_ = rootCmd.Execute()
	// Output:
	// logging level is: 4="warning"
	// logging level is: -4="debug"
}
