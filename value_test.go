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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("enum values", func() {

	Context("scalars", func() {

		It("retrieves the current enum value", func() {
			f := fmFoo
			es := enumScalar[FooModeTest]{v: &f}
			Expect(es.Get()).To(Equal(fmFoo))
		})

		DescribeTable("stringifies",
			func(e FooModeTest, expected string) {
				es := enumScalar[FooModeTest]{v: &e}
				m := newEnumMapper(FooModeIdentifiersTest, EnumCaseInsensitive)
				Expect(es.String(m)).To(Equal(expected))
			},
			Entry("fmBar", fmBar, "bar"), // sic! returns canonical name, which is "bar"
			Entry("unknown", FooModeTest(0), "<unknown>"),
		)

		It("sets a new enum value", func() {
			f := fmFoo
			es := enumScalar[FooModeTest]{v: &f}
			m := newEnumMapper(FooModeIdentifiersTest, EnumCaseInsensitive)
			Expect(es.Set("Bar", m)).To(Succeed())
			Expect(es.Get()).To(Equal(fmBar))
		})

		It("rejects setting an unknown textual representation", func() {
			f := fmFoo
			es := enumScalar[FooModeTest]{v: &f}
			m := newEnumMapper(FooModeIdentifiersTest, EnumCaseInsensitive)
			Expect(es.Set("Barumph", m)).NotTo(Succeed())
		})

		DescribeTable("completion",
			func(toc string, expected []string) {
				c := (&enumScalar[FooModeTest]{}).NewCompletor(FooModeIdentifiersTest, nil)
				actual, _ := c(nil, nil, toc)
				Expect(actual).To(ConsistOf(expected))
			},
			Entry("returns all values", "", []string{"foo", "bar", "Bar", "baz"}),
			Entry("always returns all values without filtering", "b", []string{"foo", "bar", "Bar", "baz"}),
		)

		It("completes with help", func() {
			c := (&enumScalar[FooModeTest]{}).NewCompletor(FooModeIdentifiersTest, FooModeHelp)
			actual, _ := c(nil, nil, "")
			Expect(actual).To(ConsistOf([]string{
				"foo\tfoo it",
				"bar\tbar IT!",
				"Bar\tbar IT!",
				"baz\tbaz nit!!",
			}))
		})

	})

	Context("slices", func() {

		m := newEnumMapper(FooModeIdentifiersTest, EnumCaseInsensitive)
		var es enumSlice[FooModeTest]

		BeforeEach(func() {
			sf := []FooModeTest{fmFoo, fmBar, 0}
			es = enumSlice[FooModeTest]{v: &sf}
		})

		It("retrieves the current enum value", func() {
			Expect(es.Get()).To(ConsistOf(fmFoo, fmBar, FooModeTest(0)))
		})

		It("stringifies", func() {
			Expect(es.String(m)).To(Equal("[foo,bar,<unknown>]"))
			sf := []FooModeTest{}
			es := enumSlice[FooModeTest]{v: &sf}
			Expect(es.String(m)).To(Equal("[]"))
		})

		It("sets a new enum value", func() {
			Expect(es.Set("baz", m)).To(Succeed())
			Expect(es.Get()).To(ConsistOf(fmBaz))
			Expect(es.Set("foo", m)).To(Succeed())
			Expect(es.Get()).To(ConsistOf(fmBaz, fmFoo))
			Expect(es.Set("Baz", m)).To(Succeed())
			Expect(es.Get()).To(ConsistOf(fmBaz, fmFoo))
		})

		DescribeTable("rejects setting an unknown textual representation",
			func(value string) {
				Expect(es.Set(value, m)).NotTo(Succeed())
			},
			Entry(nil, "bajazzo"),
			Entry(nil, "foo,bajazzo"),
			Entry(`""`, ""),
		)

		DescribeTable("completion",
			func(toc string, expected []string) {
				c := (&enumSlice[FooModeTest]{}).NewCompletor(FooModeIdentifiersTest, nil)
				actual, _ := c(nil, nil, toc)
				Expect(actual).To(ConsistOf(expected))
			},
			Entry("returns all values", "", []string{"foo", "bar", "Bar", "baz"}),
			Entry("always returns all values without filtering", "b",
				[]string{"foo", "bar", "Bar", "baz"}),
			Entry("returns all remaining values", "f",
				[]string{"foo", "bar", "Bar", "baz"}),
			Entry(nil, "foo,", []string{"foo,bar", "foo,Bar", "foo,baz"}),
			Entry(nil, "foo,bar,", []string{"foo,bar,Bar", "foo,bar,baz"}),
			Entry(nil, "bar,baz,f", []string{"bar,baz,foo", "bar,baz,Bar"}),
			Entry("ignores non-existing elements", "foo,koo,",
				[]string{"foo,koo,bar", "foo,koo,Bar", "foo,koo,baz"}),
		)

	})

})
