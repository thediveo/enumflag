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

// ① Defines a new enum flag type.
type FooMode enumflag.Flag

// ② Defines the enumeration values for our new FooMode enum flag type.
const (
    Foo FooMode = iota
    Bar
)

// ③ Implements the methods required by spf13/cobra in order to use the enum as
// a flag.
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

// ⑤ Now use the FooMode enum flag.
var foomode FooMode

func main() {
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
}
```

> **Important:** always define a separate type for each of your enumeration
> flag types. Behind the scenes, `enumflag` caches the enum mappings based on
> enumeration flag type.

The boilerplate pattern is always the same; unfortunately due to the Golang
language design we have to live with this boilerplate.

1. Define your own new enumeration type, such as `type FooMode enumflag.Flag`.
2. Define the constants in your enumeration.
3. Implement the methods `String` and `Set` by routing them into
   `enumflag.String` and `enumflag.Set` respectively. Make sure your `Type`
   method returns some sensible and especially unique flag type name.
4. Implement method `Enums` to return the enum value to textual representation
   mapping, as well as the desired case sensitivity.
5. Somewhere, declare a flag variable of your enum flag type.
6. Wire up your flag variable to its flag long and short names, et cetera.

## Copyright and License

`lxkns` is Copyright 2020 Harald Albrecht, and licensed under the Apache
License, Version 2.0.
