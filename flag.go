// Copyright 2020, 2022 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package enumflag

import (
	"golang.org/x/exp/constraints"
)

// Flag represents a CLI (enumeration) flag which can take on only a single
// enumeration value out of a fixed set of enumeration values. Applications
// using the enumflag package might want to derive their enumeration flags from
// Flag, such as "type MyFoo enumflag.Flag", but they don't need to. The only
// requirement for user-defined enumeration flags is that they must be
// ("somewhat") compatible with the Flag type, or more precise: user-defined
// enumerations must satisfy constraints.Integer.
type Flag uint

// EnumCaseSensitivity specifies whether the textual representations of enum
// values are considered to be case sensitive, or not.
type EnumCaseSensitivity bool

// Controls whether the textual representations for enum values are case
// sensitive, or not.
const (
	EnumCaseInsensitive EnumCaseSensitivity = false
	EnumCaseSensitive   EnumCaseSensitivity = true
)

// EnumFlagValue wraps a user-defined enum type value and implements the
// [github.com/spf13/pflag.Value] interface, so the user's enum type value can
// directly be used with the fine pflag drop-in package for Golang CLI flags.
type EnumFlagValue[E constraints.Integer] struct {
	value    enumValue[E]  // enum value of a user-defined enum scalar or slice type.
	enumtype string        // user-friendly name of the user-defined enum type.
	names    enumMapper[E] // enum value names.
}

// enumValue supports getting, setting, and stringifying an scalar or slice enum
// enumValue.
//
// Do I smell preemptive interfacing here...? Now watch the magic of cleanest
// code: by just moving the interface type from the source file with the struct
// types to the source file with the consumer we achieve immediate Go
// perfectness!
type enumValue[E constraints.Integer] interface {
	Get() any
	Set(val string, names enumMapper[E]) error
	String(names enumMapper[E]) string
}

// New wraps a given enum variable (satisfying constraints.Integer) so that it
// can be used as a flag Value with pflag.Var and pflag.VarP.
func New[E constraints.Integer](flag *E, typename string, mapping EnumIdentifiers[E], sensitivity EnumCaseSensitivity) *EnumFlagValue[E] {
	if flag == nil {
		panic("New requires flag to be a non-nil pointer to an enum value satisfying constraints.Integer")
	}
	if mapping == nil {
		panic("New requires mapping not to be nil")
	}
	return &EnumFlagValue[E]{
		value:    &enumScalar[E]{v: flag},
		enumtype: typename,
		names:    newEnumMapper(mapping, sensitivity),
	}
}

func NewSlice[E constraints.Integer](flag *[]E, typename string, mapping EnumIdentifiers[E], sensitivity EnumCaseSensitivity) *EnumFlagValue[E] {
	if flag == nil {
		panic("NewSlice requires flag to be a non-nil pointer to an enum value slice satisfying []constraints.Integer")
	}
	if mapping == nil {
		panic("NewSlice requires mapping not to be nil")
	}
	return &EnumFlagValue[E]{
		value:    &enumSlice[E]{v: flag},
		enumtype: typename,
		names:    newEnumMapper(mapping, sensitivity),
	}
}

// Set sets the enum flag to the specified enum value. If the specified value
// isn't a valid enum value, then the enum flag won't be set and an error is
// returned instead.
func (e *EnumFlagValue[E]) Set(val string) error {
	return e.value.Set(val, e.names)
}

// String returns the textual representation of an enumeration (flag) value. In
// case multiple textual representations (=identifiers) exist for the same
// enumeration value, then only the first textual representation is returned,
// which is considered to be the canonical one.
func (e *EnumFlagValue[E]) String() string { return e.value.String(e.names) }

// Type returns the name of the flag value type. The type name is used in error
// messages.
func (e *EnumFlagValue[E]) Type() string { return e.enumtype }

// Get returns the current enum value for convenience. Please note that the enum
// value is either scalar or slice, depending on how the enum flag was created.
func (e *EnumFlagValue[E]) Get() any { return e.value.Get() }
