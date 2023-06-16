// Copyright 2023 Harald Albrecht.
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

package main

import (
	"bytes"

	"github.com/spf13/cobra"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("enumflag-testing canary", func() {

	var rootCmd *cobra.Command
	var out *bytes.Buffer

	BeforeEach(func() {
		rootCmd = newRootCmd()
		out = &bytes.Buffer{}
		rootCmd.SetOutput(out)
	})

	It("has a hidden __complete command", func() {
		rootCmd.SetArgs([]string{"__complete", "t"})
		Expect(rootCmd.Execute()).To(Succeed())
		Expect(out.String()).To(MatchRegexp(`test\n:\d+\nCompletion ended with directive: .+`))
	})

	It("lists the completion command", func() {
		rootCmd.SetArgs([]string{"-h"})
		Expect(rootCmd.Execute()).To(Succeed())
		Expect(out.String()).To(MatchRegexp(`Available Commands:\n\s+completion\s+ Generate .* shell`))
	})

	It("generates a shell completion script", func() {
		rootCmd.SetArgs([]string{"completion", "bash"})
		Expect(rootCmd.Execute()).To(Succeed())
		Expect(out.String()).To(MatchRegexp(`^# bash completion V2 for`))
	})

})