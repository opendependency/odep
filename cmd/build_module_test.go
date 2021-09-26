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
	"fmt"
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/opendependency/odep/cmd"
)

var _ = Describe("Build Module Command", func() {
	const testModuleJSON = `{"namespace":"com.example.shop","name":"products","type":"go","version":{"name":"v1.1.1","schema":"org.semver.v2","replaces":["v1.1.0"]},"annotations":{"key1":"value1","key2":"value2"},"dependencies":[{"namespace":"com.example.shop","name":"web-libs","type":"go","version":"v3.2.1","direction":"UPSTREAM"},{"namespace":"com.example.shop","name":"utils","type":"go","version":"v4.1.1","direction":"UPSTREAM"},{"namespace":"com.example.shop","name":"products","type":"org.openapis","version":"v1.3.4","direction":"DOWNSTREAM"}]}`
	const testModuleJSONPretty = `{
  "namespace": "com.example.shop",
  "name": "products",
  "type": "go",
  "version": {
    "name": "v1.1.1",
    "schema": "org.semver.v2",
    "replaces": [
      "v1.1.0"
    ]
  },
  "annotations": {
    "key1": "value1",
    "key2": "value2"
  },
  "dependencies": [
    {
      "namespace": "com.example.shop",
      "name": "web-libs",
      "type": "go",
      "version": "v3.2.1",
      "direction": "UPSTREAM"
    },
    {
      "namespace": "com.example.shop",
      "name": "utils",
      "type": "go",
      "version": "v4.1.1",
      "direction": "UPSTREAM"
    },
    {
      "namespace": "com.example.shop",
      "name": "products",
      "type": "org.openapis",
      "version": "v1.3.4",
      "direction": "DOWNSTREAM"
    }
  ]
}`
	const testModuleYAMLAlphabeticSortedKeys = `annotations:
  key1: value1
  key2: value2
dependencies:
- direction: UPSTREAM
  name: web-libs
  namespace: com.example.shop
  type: go
  version: v3.2.1
- direction: UPSTREAM
  name: utils
  namespace: com.example.shop
  type: go
  version: v4.1.1
- direction: DOWNSTREAM
  name: products
  namespace: com.example.shop
  type: org.openapis
  version: v1.3.4
name: products
namespace: com.example.shop
type: go
version:
  name: v1.1.1
  replaces:
  - v1.1.0
  schema: org.semver.v2
`
	const testModuleYAMLLogicalSortedKeys = `---
namespace: com.example.shop
name: products
type: go
version:
  name: v1.1.1
  replaces:
  - v1.1.0
  schema: org.semver.v2
annotations:
  key1: value1
  key2: value2
dependencies:
- namespace: com.example.shop
  name: web-libs
  type: go
  version: v3.2.1
  direction: UPSTREAM
- namespace: com.example.shop
  name: utils
  type: go
  version: v4.1.1
  direction: UPSTREAM
- namespace: com.example.shop
  name: products
  type: org.openapis
  version: v1.3.4
  direction: DOWNSTREAM
`

	var (
		rootCmd     *cobra.Command
		rootCmdArgs []string

		stdOut *strings.Builder
		stdErr *strings.Builder
	)

	BeforeEach(func() {
		rootCmd = cmd.NewRootCommand()
		rootCmdArgs = []string{"build", "module"}

		stdOut = &strings.Builder{}
		stdErr = &strings.Builder{}

		rootCmd.SetOut(stdOut)
		rootCmd.SetErr(stdErr)
	})

	Context("command is executed", func() {
		JustBeforeEach(func() {
			rootCmd.SetArgs(rootCmdArgs)

			_ = rootCmd.Execute()
		})

		When("no flags provided", func() {
			It("should print validation error to stderr", func() {
				Expect(stdErr.String()).To(Equal("Error: validation failed: namespace: must have at least 1 characters\n"))
			})

			It("should not write to stdout", func() {
				Expect(stdOut.String()).To(Equal(""))
			})
		})

		When("module is built from file only", func() {

			When("file format is json", func() {
				BeforeEach(func() {
					f, err := os.CreateTemp("", "module*.json")
					if err != nil {
						Fail(fmt.Sprintf("could not create temporary module file: %v", err))
					}

					_, err = f.WriteString(testModuleJSON)
					if err != nil {
						Fail(fmt.Sprintf("could not write to temporary module file: %v", err))
					}

					rootCmdArgs = append(rootCmdArgs, "-f", f.Name())
				})

				It("should print module built", func() {
					Expect(stdOut.String()).To(Equal("Module com.example.shop products go v1.1.1 built.\n"))
				})

				It("should not write to stderr", func() {
					Expect(stdErr.String()).To(Equal(""))
				})

				When("flag output is set to json", func() {
					BeforeEach(func() {
						rootCmdArgs = append(rootCmdArgs, "--output", "json")
					})

					It("should print module json to stdout", func() {
						Expect(stdOut.String()).To(Equal(testModuleJSON))
					})

					It("should not write to stderr", func() {
						Expect(stdErr.String()).To(Equal(""))
					})

					When("flag pretty is set to true", func() {
						BeforeEach(func() {
							rootCmdArgs = append(rootCmdArgs, "--pretty")
						})

						It("should print module multi-line json with indents to stdout", func() {
							Expect(stdOut.String()).To(Equal(testModuleJSONPretty))
						})

						It("should not write to stderr", func() {
							Expect(stdErr.String()).To(Equal(""))
						})
					})
				})

				When("flag output is set to yaml", func() {
					BeforeEach(func() {
						rootCmdArgs = append(rootCmdArgs, "--output", "yaml")
					})

					It("should print module yaml to stdout", func() {
						Expect(stdOut.String()).To(Equal(testModuleYAMLAlphabeticSortedKeys))
					})

					It("should not write to stderr", func() {
						Expect(stdErr.String()).To(Equal(""))
					})

					When("flag pretty is set to true", func() {
						BeforeEach(func() {
							rootCmdArgs = append(rootCmdArgs, "--pretty")
						})

						It("should print module yaml to stdout", func() {
							Expect(stdOut.String()).To(Equal(testModuleYAMLAlphabeticSortedKeys))
						})

						It("should not write to stderr", func() {
							Expect(stdErr.String()).To(Equal(""))
						})
					})
				})
			})

			When("file format is yaml", func() {

				BeforeEach(func() {
					f, err := os.CreateTemp("", "module*.yaml")
					if err != nil {
						Fail(fmt.Sprintf("could not create temporary module file: %v", err))
					}

					_, err = f.WriteString(testModuleYAMLLogicalSortedKeys)
					if err != nil {
						Fail(fmt.Sprintf("could not write to temporary module file: %v", err))
					}

					rootCmdArgs = append(rootCmdArgs, "-f", f.Name())
				})

				It("should print module built", func() {
					Expect(stdOut.String()).To(Equal("Module com.example.shop products go v1.1.1 built.\n"))
				})

				It("should not write to stderr", func() {
					Expect(stdErr.String()).To(Equal(""))
				})

				When("flag output is set to json", func() {
					BeforeEach(func() {
						rootCmdArgs = append(rootCmdArgs, "--output", "json")
					})

					It("should print module json to stdout", func() {
						Expect(stdOut.String()).To(Equal(testModuleJSON))
					})

					It("should not write to stderr", func() {
						Expect(stdErr.String()).To(Equal(""))
					})
				})

				When("flag output is set to yaml", func() {
					BeforeEach(func() {
						rootCmdArgs = append(rootCmdArgs, "--output", "yaml")
					})

					It("should print module json to stdout", func() {
						Expect(stdOut.String()).To(Equal(testModuleYAMLAlphabeticSortedKeys))
					})

					It("should not write to stderr", func() {
						Expect(stdErr.String()).To(Equal(""))
					})
				})
			})
		})

		When("module is built by flags only", func() {
			BeforeEach(func() {
				rootCmdArgs = append(rootCmdArgs,
					"--namespace", "com.example.shop",
					"--name", "products",
					"--type", "go",
					"--version-name", "v1.1.1",
					"--version-schema", "org.semver.v2",
					"--version-replaces", "v1.1.0",
					"--annotations", "key1=value1,key2=value2",
					"--upstream-dependencies", "com.example.shop:web-libs:go:v3.2.1,com.example.shop:utils:go:v4.1.1",
					"--downstream-dependencies", "com.example.shop:products:org.openapis:v1.3.4",
				)
			})

			It("should print module built", func() {
				Expect(stdOut.String()).To(Equal("Module com.example.shop products go v1.1.1 built.\n"))
			})

			It("should not write to stderr", func() {
				Expect(stdErr.String()).To(Equal(""))
			})

			When("flag output is set to json", func() {
				BeforeEach(func() {
					rootCmdArgs = append(rootCmdArgs, "--output", "json")
				})

				It("should print module json to stdout", func() {
					Expect(stdOut.String()).To(Equal(testModuleJSON))
				})

				It("should not write to stderr", func() {
					Expect(stdErr.String()).To(Equal(""))
				})
			})

			When("flag output is set to yaml", func() {
				BeforeEach(func() {
					rootCmdArgs = append(rootCmdArgs, "--output", "yaml")
				})

				It("should print module json to stdout", func() {
					Expect(stdOut.String()).To(Equal(testModuleYAMLAlphabeticSortedKeys))
				})

				It("should not write to stderr", func() {
					Expect(stdErr.String()).To(Equal(""))
				})
			})
		})
	})
})
