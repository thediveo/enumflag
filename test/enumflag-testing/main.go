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
	"os"

	"github.com/spf13/cobra"
)

const Name = "enumflag-testing"

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: Name,
		Run: func(*cobra.Command, []string) {},
	}
	testCmd := &cobra.Command{
		Use: "test",
		Run: func(*cobra.Command, []string) {},
	}
	rootCmd.AddCommand(testCmd)
	return rootCmd
}

func main() {
	// Cobra automatically adds a "__complete" command to our root command
	// behind the scenes, unless we specify one explicitly. It also adds a
	// "complete" sub command if we're adding at least one sub command.
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
