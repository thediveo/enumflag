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

func (f *FooModeTest) String() string     { return String(f) }
func (f *FooModeTest) Set(s string) error { return Set(f, s) }
func (f *FooModeTest) Enums() (interface{}, EnumCaseSensitivity) {
	return FooModeTestIdentifiers, EnumCaseSensitive
}

// Enumeration constants/values.
const (
	fmFoo FooModeTest = iota
	fmBar
	fmBaz
)

// Enumeration identifiers mapped to their corresponding constants.
var FooModeTestIdentifiers = map[FooModeTest][]string{
	fmFoo: {"foo"},
	fmBar: {"bar", "Bar"},
	fmBaz: {"baz"},
}

// Another new enumeration type.
type BarzModeTest Flag

func (f *BarzModeTest) String() string     { return String(f) }
func (f *BarzModeTest) Set(s string) error { return Set(f, s) }
func (f *BarzModeTest) Enums() (interface{}, EnumCaseSensitivity) {
	return BarzModeTestIdentifiers, EnumCaseSensitive
}

const (
	bmbarz BarzModeTest = iota
	bmBarz
	bmBARZ
)

var BarzModeTestIdentifiers = map[BarzModeTest][]string{
	bmbarz: {"barz"},
	bmBarz: {"Barz"},
	bmBARZ: {"BARZ"},
}

type Wrong Flag

func (w *Wrong) String() string     { return String(w) }
func (w *Wrong) Set(s string) error { return Set(w, s) }
func (w *Wrong) Enums() (interface{}, EnumCaseSensitivity) {
	return "foobar", EnumCaseSensitive
}

var _ = Describe("flag", func() {

	It("returns the canonical textual representation", func() {
		foomode := fmBar
		Expect(foomode.String()).To(Equal(FooModeTestIdentifiers[fmBar][0]))
	})

	It("denies setting invalid values", func() {
		var foomode FooModeTest
		err := foomode.Set("FOOBAR")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("must be 'bar'/'Bar', 'baz', 'foo'"))
	})

	It("sets the enumeration value from text", func() {
		var foomode FooModeTest
		var barzmode BarzModeTest

		Expect(foomode.Set("foo")).NotTo(HaveOccurred())
		Expect(barzmode.Set("BARZ")).NotTo(HaveOccurred())
		Expect(foomode.Set("Bar")).NotTo(HaveOccurred())

		Expect(foomode).To(Equal(fmBar))
		Expect(barzmode).To(Equal(bmBARZ))
	})

	It("accepts only uint/int enum flags", func() {
		var sf string
		Expect(func() { _ = Set(sf, "6.66") }).To(Panic())
	})

	It("checks for a map", func() {
		var w Wrong
		Expect(func() { _ = w.Set("abc") }).To(Panic())
	})

})
