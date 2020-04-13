# CLI Enumeration Flags

[![GoDoc](https://godoc.org/github.com/thediveo/enumflag?status.svg)](http://godoc.org/github.com/thediveo/enumflag)
[![GitHub](https://img.shields.io/github/license/thediveo/enumflag)](https://img.shields.io/github/license/thediveo/enumflag)
![build and test](https://github.com/thediveo/enumflag/workflows/build%20and%20test/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/thediveo/enumflag)](https://goreportcard.com/report/github.com/thediveo/enumflag)
![Coverage](https://img.shields.io/badge/coverage-135.72%25-darkcyan)

`enumflag` is a Golang package which supplements the Golang CLI flag packages
[spf13/cobra](https://github.com/spf13/cobra) and
[spf13/pflag](https://github.com/spf13/pflag) with enumeration flags, including
support for enumeration slices.

For instance, users can specify enum flags as `--mode=foo` or `--mode=bar`,
where `foo` and `bar` are valid enumeration values. Other values which are not
part of the set of allowed enumeration values cannot be set and raise CLI flag
errors. In case of an enumeration _slice_ flag users can specify multiple
enumeration values either with a single flag `--mode=foo,bar` or multiple flag
calls, such as `--mode=foo --mode=bar`.

Application programmers then simply deal with enumeration values in form of
uints (or ints), liberated from parsing strings and validating enumeration
flags.

## Alternatives

In case you are just interested in string-based one-of-a-set flags, then the
following packages offer you a minimalist approach:

- [hashicorp/packer/helper/enumflag](https://godoc.org/github.com/hashicorp/packer/helper/enumflag)
  really is a reduced-to-the-max version without any whistles and bells.

- [creachadair/goflags/enumflag](https://godoc.org/github.com/creachadair/goflags/enumflag)
  has a similar, but slightly more elaborate API with additional "indices" for
  enumeration values.

But if you instead want to handle one-of-a-set flags as properly typed
enumerations instead of strings, or if you need (multiple-of-a-set) slice
support, then please read on.

## How To Use: Properly Typed Enum Flag

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

## How To Use: Slice of Enums

For a slice of enumerations, simply declare your variable to be a slice of your
enumeration type and then use `enumflag.NewSlice(...)` instead of
`enumflag.New(...)`.

```go
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
    // ⑤ Define the CLI flag parameters for your wrapped enumm slice flag.
    rootCmd.PersistentFlags().VarP(
        enumflag.NewSlice(&moomode, "mode", MooModeIds, enumflag.EnumCaseInsensitive),
        "mode", "m",
        "can be any combination of 'moo', 'møø', 'mimimi'")

    rootCmd.SetArgs([]string{"--mode", "Moo,møø"})
    _ = rootCmd.Execute()
}
```

## Copyright and License

`lxkns` is Copyright 2020 Harald Albrecht, and licensed under the Apache
License, Version 2.0.
