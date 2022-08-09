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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/thediveo/enumflag"
)

var _ = Describe("flag slices", func() {

	It("returns the canonical textual representation", func() {
		foomodes := []FooModeTest{fmFoo, fmBar}
		val := enumflag.NewSlice(&foomodes, "mode", FooModeIdentifiersTest, enumflag.EnumCaseInsensitive)
		Expect(val.String()).To(Equal("[foo,bar]"))
		Expect(val.Type()).To(Equal("mode"))
	})

	It("replaces, then merges", func() {
		foomodes := []FooModeTest{fmFoo, fmBar}
		val := enumflag.NewSlice(&foomodes, "mode", FooModeIdentifiersTest, enumflag.EnumCaseInsensitive)
		Expect(val.Set("baz")).NotTo(HaveOccurred())
		Expect(val.String()).To(Equal("[baz]"))

		Expect(val.Set("baz")).NotTo(HaveOccurred())
		Expect(val.String()).To(Equal("[baz]"))

		Expect(val.Set("foo")).NotTo(HaveOccurred())
		Expect(val.String()).To(Equal("[baz,foo]"))

		Expect(val.Set("foo,bar,baz")).NotTo(HaveOccurred())
		Expect(val.String()).To(Equal("[baz,foo,bar]"))
	})

	It("accepts only uint/int enum slice flags by reference", func() {
		var sf string
		Expect(func() {
			_ = enumflag.NewSlice(&sf, "string", FooModeIdentifiersTest, enumflag.EnumCaseSensitive)
		}).To(Panic())
		var foomodes []FooModeTest
		Expect(func() {
			_ = enumflag.NewSlice(foomodes, "mode", "abc", enumflag.EnumCaseSensitive)
		}).To(Panic())
	})

	It("rejects setting invalid values", func() {
		var foomodes []FooModeTest
		val := enumflag.NewSlice(&foomodes, "mode", FooModeIdentifiersTest, enumflag.EnumCaseSensitive)
		err := val.Set("FOOBAR")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("must be 'bar'/'Bar', 'baz', 'foo'"))
	})

})
