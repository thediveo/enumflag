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
	"github.com/spf13/cobra"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// EnumHelp maps enumeration values to their corresponding help texts. These
// help texts should contain just the help text but without any "foo\t" enum
// value name prefix. The reason is that enumflag will automatically register
// the correct completion text. Please note that it isn't necessary to supply
// any help texts in order to register enum flag completion.
type EnumHelp[E constraints.Integer] map[E]string

// Completor tells cobra how to complete a flag. See also [dynamic flag
// completion].
//
// [dynamic flag completion]: https://github.com/spf13/cobra/blob/main/shell_completions.md#specify-dynamic-flag-completion
type Completor func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)

// newCompletor returns a function that can be registered with cobra flags in
// order to provide enum flag name completion.
func newCompletor[E constraints.Integer](enums EnumIdentifiers[E], help EnumHelp[E]) Completor {
	completions := []string{}
	for enumval, enumnames := range enums {
		helptext := ""
		if text, ok := help[enumval]; ok {
			helptext = "\t" + text
		}
		// complete not only the canonical enum value name, but also all other
		// (alias) names.
		for _, name := range enumnames {
			completions = append(completions, name+helptext)
		}
	}
	slices.Sort(completions)
	return func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return completions, cobra.ShellCompDirectiveDefault
	}
}
