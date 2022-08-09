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

var _ = Describe("enum name-value mapper", func() {

	DescribeTable("looks up name for value",
		func(value FooModeTest, expectedNames []string) {
			mapper := newEnumMapper(FooModeIdentifiersTest, EnumCaseSensitive)
			Expect(mapper.Lookup(value)).To(Equal(expectedNames))
		},
		Entry("fmBar", fmBar, FooModeIdentifiersTest[fmBar]),
		Entry("fool", FooModeTest(0), nil),
	)

	DescribeTable("looks up value for name",
		func(name string, sensitivity EnumCaseSensitivity, expectedValue FooModeTest) {
			mapper := newEnumMapper(FooModeIdentifiersTest, sensitivity)
			Expect(mapper.ValueOf(name)).To(Equal(expectedValue))
		},
		Entry("baz", "baz", EnumCaseSensitive, fmBaz),
		Entry("Baz/i", "Baz", EnumCaseInsensitive, fmBaz),
		Entry("bar", "bar", EnumCaseSensitive, fmBar),
		Entry("Bar", "Bar", EnumCaseSensitive, fmBar),
	)

	DescribeTable("returns helpful error when lookup fails",
		func(name string, sensitivity EnumCaseSensitivity) {
			mapper := newEnumMapper(FooModeIdentifiersTest, sensitivity)
			Expect(mapper.ValueOf(name)).Error().To(MatchError("must be 'bar'/'Bar', 'baz', 'foo'"))
		},
		Entry("fool", "fool", EnumCaseSensitive),
		Entry("fool/i", "fool", EnumCaseInsensitive),
		Entry("BAr", "BAr", EnumCaseSensitive),
	)

})
