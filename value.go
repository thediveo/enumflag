// Copyright 2022 Harald Albrecht.
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
	"strings"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// unknown is the textual representation of an unknown enum value, that is, when
// the enum value to name mapping doesn't have any idea about a particular enum
// value.
const unknown = "<unknown>"

// enumScalar represents a mutable, single enumeration value that can be
// retrieved, set, and stringified.
type enumScalar[E constraints.Integer] struct {
	v *E
}

// Get returns the scalar enum value.
func (s *enumScalar[E]) Get() any { return *s.v }

// Set the value to the new scalar enum value corresponding to the passed
// textual representation, using the additionally specified text-to-value
// mapping. If the specified textual representation doesn't match any of the
// defined ones, an error is returned instead and the value isn't changed.
func (s *enumScalar[E]) Set(val string, names enumMapper[E]) error {
	enumcode, err := names.ValueOf(val)
	if err != nil {
		return err
	}
	*s.v = enumcode
	return nil
}

// String returns the textual representation of the scalar enum value, using the
// specified text-to-value mapping.
func (s *enumScalar[E]) String(names enumMapper[E]) string {
	if ids := names.Lookup(*s.v); len(ids) > 0 {
		return ids[0]
	}
	return unknown
}

// enumSlice represents a slice of enumeration values that can be retrieved,
// set, and stringified.
type enumSlice[E constraints.Integer] struct {
	v     *[]E
	merge bool // replace the complete slice or merge values?
}

// Get returns the slice enum values.
func (s *enumSlice[E]) Get() any { return *s.v }

// Set or merge one or more values of the new scalar enum value corresponding to
// the passed textual representation, using the additionally specified
// text-to-value mapping. If the specified textual representation doesn't match
// any of the defined ones, an error is returned instead and the value isn't
// changed. The first call to Set will always clear any previous default value.
// All subsequent calls to Set will merge the specified enum values with the
// current enum values.
func (s *enumSlice[E]) Set(val string, names enumMapper[E]) error {
	// First parse and convert the textual enum values into their
	// program-internal codes.
	ids := strings.Split(val, ",")
	enumvals := make([]E, 0, len(ids)) // ...educated guess
	for _, id := range ids {
		enumval, err := names.ValueOf(id)
		if err != nil {
			return err
		}
		enumvals = append(enumvals, enumval)
	}
	if !s.merge {
		// Replace any existing default enum value set on first Set().
		*s.v = enumvals
		s.merge = true // ...and next time: merge.
		return nil
	}
	// Later, merge with the existing enum values.
	for _, enumval := range enumvals {
		if slices.Index(*s.v, enumval) >= 0 {
			continue
		}
		*s.v = append(*s.v, enumval)
	}
	return nil
}

// String returns the textual representation of the slice enum value, using the
// specified text-to-value mapping.
func (s *enumSlice[E]) String(names enumMapper[E]) string {
	n := make([]string, 0, len(*s.v))
	for _, enumval := range *s.v {
		if enumnames := names.Lookup(enumval); len(enumnames) > 0 {
			n = append(n, enumnames[0])
			continue
		}
		n = append(n, unknown)
	}
	return "[" + strings.Join(n, ",") + "]"
}
