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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("flag", func() {

	Context("scalar enum flag", func() {

		It("returns the canonical textual representation", func() {
			foomode := fmBar
			val := New(&foomode, "mode", FooModeIdentifiersTest, EnumCaseInsensitive)
			Expect(val.String()).To(Equal("bar"))
			Expect(val.Type()).To(Equal("mode"))
		})

		It("rejects setting invalid values", func() {
			var foomode FooModeTest
			val := New(&foomode, "mode", FooModeIdentifiersTest, EnumCaseSensitive)
			Expect(val.Set("FOOBAR")).To(MatchError("must be 'bar'/'Bar', 'baz', 'foo'"))
		})

		It("sets the enumeration value from text", func() {
			var foomode FooModeTest
			val := New(&foomode, "mode", FooModeIdentifiersTest, EnumCaseSensitive)

			Expect(val.Set("foo")).NotTo(HaveOccurred())
			Expect(val.Set("Bar")).NotTo(HaveOccurred())
			Expect(foomode).To(Equal(fmBar))
			Expect(val.Get()).To(Equal(fmBar))
		})

	})

	Context("slice enum flag", func() {

		It("returns the canonical textual representation", func() {
			foomodes := []FooModeTest{fmBar, fmFoo}
			val := NewSlice(&foomodes, "modes", FooModeIdentifiersTest, EnumCaseInsensitive)
			Expect(val.String()).To(Equal("[bar,foo]"))
			Expect(val.Type()).To(Equal("modes"))
		})

	})

	When("passing nil", func() {

		It("panics", func() {
			Expect(func() {
				_ = New[FooModeTest](nil, "foo", nil, EnumCaseInsensitive)
			}).To(PanicWith(MatchRegexp("New requires flag to be a non-nil pointer to .*")))
			Expect(func() {
				var f FooModeTest
				_ = New(&f, "foo", nil, EnumCaseInsensitive)
			}).To(PanicWith(MatchRegexp("New requires mapping not to be nil")))

			Expect(func() {
				_ = NewSlice[FooModeTest](nil, "foo", nil, EnumCaseInsensitive)
			}).To(PanicWith(MatchRegexp("NewSlice requires flag to be a non-nil pointer to .*")))
			Expect(func() {
				var f []FooModeTest
				_ = NewSlice(&f, "foo", nil, EnumCaseInsensitive)
			}).To(PanicWith(MatchRegexp("NewSlice requires mapping not to be nil")))
		})

	})

})
