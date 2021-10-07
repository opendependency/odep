/*
Copyright © 2021 The OpenDependency Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/opendependency/odep/cmd"
)

var _ = Describe("Root Command", func() {
	var (
		rootCmd *cobra.Command

		stdOut *strings.Builder
		stdErr *strings.Builder
	)

	BeforeEach(func() {
		rootCmd = cmd.NewRootCommand(cmd.NewCommandContext(nil))

		stdOut = &strings.Builder{}
		stdErr = &strings.Builder{}

		rootCmd.SetOut(stdOut)
		rootCmd.SetErr(stdErr)
	})

	Context("command is executed", func() {
		JustBeforeEach(func() {
			_ = rootCmd.Execute()
		})

		When("no sub-command is called", func() {
			It("should print help to stdout", func() {
				Expect(stdOut.String()).To(Equal(`odep manages OpenDependency modules.

Usage:
  odep [command]

Available Commands:
  build       Builds OpenDependency artifacts.
  completion  generate the autocompletion script for the specified shell
  help        Help about any command

Flags:
  -h, --help   help for odep

Use "odep [command] --help" for more information about a command.
`))
			})

			It("should not write to stderr", func() {
				Expect(stdErr.String()).To(Equal(""))
			})
		})
	})
})
