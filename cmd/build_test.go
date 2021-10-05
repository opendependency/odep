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

var _ = Describe("Build Command", func() {
	var (
		rootCmd     *cobra.Command
		rootCmdArgs []string

		stdOut *strings.Builder
		stdErr *strings.Builder
	)

	BeforeEach(func() {
		rootCmd = cmd.NewRootCommand(cmd.NewContext(nil))
		rootCmdArgs = []string{"build"}

		stdOut = &strings.Builder{}
		stdErr = &strings.Builder{}

		rootCmd.SetOut(stdOut)
		rootCmd.SetErr(stdErr)
	})

	Context("command is executed", func() {
		JustBeforeEach(func() {
			rootCmd.SetArgs(rootCmdArgs)

			err := rootCmd.Execute()
			Expect(err).To(BeNil())
		})

		When("no sub-command is called", func() {
			It("should print help to stdout", func() {
				Expect(stdOut.String()).To(Equal(`Builds OpenDependency artifacts.

Usage:
  odep build [command]

Available Commands:
  module      Builds a module.

Flags:
  -h, --help   help for build

Use "odep build [command] --help" for more information about a command.
`))
			})

			It("should not write to stderr", func() {
				Expect(stdErr.String()).To(Equal(""))
			})
		})
	})
})
