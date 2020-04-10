// Copyright 2020 Harald Albrecht.
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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Our new enumeration type.
type FooModeTest Flag

func (fm *FooModeTest) String() string     { return String(fm) }
func (fm *FooModeTest) Set(s string) error { return Set(fm, s) }
func (fm *FooModeTest) Enums() (interface{}, EnumCaseSensitivity) {
	return FooModeIdentifiers, EnumCaseSensitive
}

// Enumeration constants/values.
const (
	Foo FooModeTest = iota
	Bar
	Baz
)

// Enumeration identifiers mapped to their corresponding constants.
var FooModeIdentifiers = map[FooModeTest][]string{
	Foo: {"foo"},
	Bar: {"bar", "Bar"},
	Baz: {"baz"},
}

var _ = Describe("flag", func() {

	It("returns the canonical textual representation", func() {
		foomode := Bar
		Expect(foomode.String()).To(Equal(FooModeIdentifiers[Bar][0]))
	})

	It("denies setting invalid values", func() {
		var foomode FooModeTest
		err := foomode.Set("FOOBAR")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("must be 'bar'/'Bar', 'baz', 'foo'"))
	})

	It("sets the enumeration value from text", func() {
		var foomode FooModeTest
		Expect(foomode.Set("Bar")).NotTo(HaveOccurred())
		Expect(foomode).To(Equal(Bar))
	})

})
