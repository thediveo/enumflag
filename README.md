# CLI Enumeration Flags

[![GoDoc](https://godoc.org/github.com/thediveo/enumflag?status.svg)](http://godoc.org/github.com/thediveo/enumflag)
[![GitHub](https://img.shields.io/github/license/thediveo/enumflag)](https://img.shields.io/github/license/thediveo/enumflag)
![build and test](https://github.com/thediveo/enumflag/workflows/build%20and%20test/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/thediveo/enumflag)](https://goreportcard.com/report/github.com/thediveo/enumflag)

`enumflag` is a Golang package which supplements the Golang CLI flag handling
packages [spf13/cobra](https://github.com/spf13/cobra) and
[spf13/pflag](https://github.com/spf13/pflag) with enumeration flags.

## How To Use

Without further ado, here's how to define and use enum flags in your own
applications:

```go
import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/thediveo/enumflag"
)

// Defines a new enum flag type.
type FooMode enumflag.Flag

// Defines the enumeration values for our new FooMode enum flag type.
const (
    Foo FooMode = iota
    Bar
)

// Implements the methods required by spf13/cobra in order to use the enum as
// a flag.
func (f *FooMode) String() string     { return enumflag.String(f) }
func (f *FooMode) Set(s string) error { return enumflag.Set(f, s) }
func (f *FooMode) Type() string       { return "foomode" }

// Implements the method required by enumflag to map enum values to their
// textual identifiers.
func (f *FooMode) Enums() (interface{}, enumflag.EnumCaseSensitivity) {
    return map[FooMode][]string{
        Foo: {"foo"},
        Bar: {"bar"},
    }, enumflag.EnumCaseInsensitive
}

// Now use the FooMode enum flag...
var foomode FooMode

func main() {
    rootCmd := &cobra.Command{
        Run: func(_ *cobra.Command, _ []string) {
            fmt.Printf("mode is: %d=%q\n", foomode, foomode.String())
        },
    }
    rootCmd.PersistentFlags().VarP(
        &foomode,
        "mode", "m",
        "foos the output; can be 'foo' or 'bar'")
    rootCmd.SetArgs([]string{"--mode", "bAr"})
    _ = rootCmd.Execute()
}
```

## Copyright and License

`lxkns` is Copyright 2020 Harald Albrecht, and licensed under the Apache
License, Version 2.0.
