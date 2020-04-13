# CLI Enumeration Flags

[![GoDoc](https://godoc.org/github.com/thediveo/enumflag?status.svg)](http://godoc.org/github.com/thediveo/enumflag)
[![GitHub](https://img.shields.io/github/license/thediveo/enumflag)](https://img.shields.io/github/license/thediveo/enumflag)
![build and test](https://github.com/thediveo/enumflag/workflows/build%20and%20test/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/thediveo/enumflag)](https://goreportcard.com/report/github.com/thediveo/enumflag)
![Coverage](https://img.shields.io/badge/coverage-135.72%25-darkcyan)

`enumflag` is a Golang package which supplements the Golang CLI flag handling
packages [spf13/cobra](https://github.com/spf13/cobra) and
[spf13/pflag](https://github.com/spf13/pflag) with enumeration flags.

For instance, users can specify enum flags as `--mode=foo` or `--mode=bar`,
where `foo` and `bar` are valid enumeration values. Other values which are not
part of the set of allowed enumeration values cannot be set and raise CLI flag
errors.

Application programmers then simply deal with enumeration values in form of
uints (or ints), liberated from parsing strings and validating enumeration
flags.

## How To Use

Without further ado, here's how to define and use enum flags in your own
applications...

```go
import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/thediveo/enumflag"
)

// ① Define your new enum flag type. It can be derived from enumflag.Flag, but
// it doesn't need to be as long as it is compatible with enumflag.Flag, so
// either an int or uint.
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

// ④ Now use the FooMode enum flag.
var foomode FooMode

func main() {
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
}
```

The boilerplate pattern is always the same:

1. Define your own new enumeration type, such as `type FooMode enumflag.Flag`.
2. Define the constants in your enumeration.
3. Define the mapping of the constants onto enum values (textual
   representations).
4. Somewhere, declare a flag variable of your enum flag type.
5. Wire up your flag variable to its flag long and short names, et cetera.

## Copyright and License

`lxkns` is Copyright 2020 Harald Albrecht, and licensed under the Apache
License, Version 2.0.
