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
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Our new enumeration type.
type FooModeTest Flag

// Enumeration constants/values.
const (
	fmFoo FooModeTest = iota + 1
	fmBar
	fmBaz
)

// Enumeration identifiers mapped to their corresponding constants.
var FooModeIdentifiersTest = map[FooModeTest][]string{
	fmFoo: {"foo"},
	fmBar: {"bar", "Bar"},
	fmBaz: {"baz"},
}

func TestEnumFlag(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "enumflag")
}
