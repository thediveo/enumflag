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

package enumflag_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thediveo/enumflag"
)

// Our new enumeration type.
type FooModeTest enumflag.Flag

// Enumeration constants/values.
const (
	fmFoo FooModeTest = iota
	fmBar
	fmBaz
)

// Enumeration identifiers mapped to their corresponding constants.
var FooModeIdentifiersTest = map[FooModeTest][]string{
	fmFoo: {"foo"},
	fmBar: {"bar", "Bar"},
	fmBaz: {"baz"},
}

var _ = Describe("flag", func() {

	It("returns the canonical textual representation", func() {
		foomode := fmBar
		val := enumflag.New(&foomode, "mode", FooModeIdentifiersTest, enumflag.EnumCaseInsensitive)
		Expect(val.String()).To(Equal(FooModeIdentifiersTest[fmBar][0]))
		Expect(val.Type()).To(Equal("mode"))
	})

	It("denies setting invalid values", func() {
		var foomode FooModeTest
		val := enumflag.New(&foomode, "mode", FooModeIdentifiersTest, enumflag.EnumCaseSensitive)
		err := val.Set("FOOBAR")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("must be 'bar'/'Bar', 'baz', 'foo'"))
	})

	It("sets the enumeration value from text", func() {
		var foomode FooModeTest
		val := enumflag.New(&foomode, "mode", FooModeIdentifiersTest, enumflag.EnumCaseSensitive)

		Expect(val.Set("foo")).NotTo(HaveOccurred())
		Expect(val.Set("Bar")).NotTo(HaveOccurred())
		Expect(foomode).To(Equal(fmBar))
		Expect(*val.Get().(*FooModeTest)).To(Equal(fmBar))
	})

	It("accepts only uint/int enum flags by reference", func() {
		var sf string
		Expect(func() {
			_ = enumflag.New(&sf, "string", FooModeIdentifiersTest, enumflag.EnumCaseSensitive)
		}).To(Panic())
		var foomode FooModeTest
		Expect(func() {
			_ = enumflag.New(foomode, "mode", "abc", enumflag.EnumCaseSensitive)
		}).To(Panic())
	})

	It("checks for a compatible map", func() {
		var foomode FooModeTest
		Expect(func() {
			_ = enumflag.New(&foomode, "mode", "abc", enumflag.EnumCaseInsensitive)
		}).To(Panic())
		ids := map[string][]string{}
		Expect(func() {
			_ = enumflag.New(&foomode, "mode", ids, enumflag.EnumCaseInsensitive)
		}).To(Panic())
	})

	It("returns <unknown> if enum value is unknown", func() {
		foomode := FooModeTest(42)
		val := enumflag.New(&foomode, "mode", FooModeIdentifiersTest, enumflag.EnumCaseInsensitive)
		Expect(val.String()).To(Equal("<unknown>"))
	})

})
