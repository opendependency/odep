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

package repository

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	spec "github.com/opendependency/go-spec/pkg/spec/v1"
	"google.golang.org/protobuf/proto"
)

var _ = Describe("file repository", func() {
	var (
		tempDir string
		repo    *fileRepository
	)

	BeforeEach(func() {
		var err error

		tempDir, err = ioutil.TempDir(os.TempDir(), "file-repository")
		if err != nil {
			Fail(err.Error())
		}

		repo, err = NewFileRepository(tempDir)
		if err != nil {
			Fail(err.Error())
		}
	})

	AfterEach(func() {
		if err := os.RemoveAll(tempDir); err != nil {
			Fail(err.Error())
		}
	})

	Context("add module", func() {

		var (
			module *spec.Module
		)

		BeforeEach(func() {
			module = nil
		})

		When("given module is nil", func() {
			BeforeEach(func() {
				module = nil
			})

			It("returns an error", func() {
				err := repo.AddModule(module)
				Expect(err).To(MatchError("module must not be nil"))
			})
		})

		When("given module does not fulfil specification", func() {
			BeforeEach(func() {
				module = &spec.Module{}
			})

			It("returns an error", func() {
				err := repo.AddModule(module)
				Expect(err).To(MatchError("module validation failed: namespace: must have at least 1 characters"))
			})
		})

		When("given module does not fulfil specification", func() {
			BeforeEach(func() {
				module = &spec.Module{}
			})

			It("returns an error", func() {
				err := repo.AddModule(module)
				Expect(err).To(MatchError("module validation failed: namespace: must have at least 1 characters"))
			})
		})

		When("given module fulfils specification", func() {
			BeforeEach(func() {
				module = &spec.Module{
					Namespace: "com.example",
					Name:      "product",
					Type:      "go",
					Version: &spec.ModuleVersion{
						Name: "v1.0.0",
					},
				}
			})

			It("returns no error", func() {
				err := repo.AddModule(module)
				Expect(err).To(BeNil())
			})
		})
	})

	Context("delete namespace", func() {

		BeforeEach(func() {
			module := &spec.Module{
				Namespace: "com.example",
				Name:      "product",
				Type:      "go",
				Version: &spec.ModuleVersion{
					Name: "v1.0.0",
				},
			}

			Expect(repo.AddModule(module)).To(BeNil())
		})

		When("given namespace is empty", func() {
			It("returns no error", func() {
				err := repo.DeleteNamespace("")
				Expect(err).To(BeNil())
			})
		})

		When("given namespace does not exist", func() {
			It("returns no error", func() {
				err := repo.DeleteNamespace("com.other")
				Expect(err).To(BeNil())
			})
		})

		When("given namespace does exist", func() {
			It("returns no error", func() {
				err := repo.DeleteNamespace("com.example")
				Expect(err).To(BeNil())
			})
		})
	})

	Context("delete module", func() {

		BeforeEach(func() {
			module := &spec.Module{
				Namespace: "com.example",
				Name:      "product",
				Type:      "go",
				Version: &spec.ModuleVersion{
					Name: "v1.0.0",
				},
			}

			Expect(repo.AddModule(module)).To(BeNil())
		})

		When("given module is empty", func() {
			It("returns no error", func() {
				err := repo.DeleteModule("com.example", "")
				Expect(err).To(BeNil())
			})
		})

		When("given module does not exist", func() {
			It("returns no error", func() {
				err := repo.DeleteModule("com.example", "unknown")
				Expect(err).To(BeNil())
			})
		})

		When("given module does exist", func() {
			It("returns no error", func() {
				err := repo.DeleteModule("com.example", "product")
				Expect(err).To(BeNil())
			})
		})
	})

	Context("delete module type", func() {

		BeforeEach(func() {
			module := &spec.Module{
				Namespace: "com.example",
				Name:      "product",
				Type:      "go",
				Version: &spec.ModuleVersion{
					Name: "v1.0.0",
				},
			}

			Expect(repo.AddModule(module)).To(BeNil())
		})

		When("given module type is empty", func() {
			It("returns no error", func() {
				err := repo.DeleteModuleType("com.example", "product", "")
				Expect(err).To(BeNil())
			})
		})

		When("given module type  does not exist", func() {
			It("returns no error", func() {
				err := repo.DeleteModuleType("com.example", "product", "unknown")
				Expect(err).To(BeNil())
			})
		})

		When("given module type does exist", func() {
			It("returns no error", func() {
				err := repo.DeleteModuleType("com.example", "product", "go")
				Expect(err).To(BeNil())
			})
		})
	})

	Context("delete module version", func() {

		BeforeEach(func() {
			module := &spec.Module{
				Namespace: "com.example",
				Name:      "product",
				Type:      "go",
				Version: &spec.ModuleVersion{
					Name: "v1.0.0",
				},
			}

			Expect(repo.AddModule(module)).To(BeNil())
		})

		When("given module version is empty", func() {
			It("returns no error", func() {
				err := repo.DeleteModuleVersion("com.example", "product", "go", "")
				Expect(err).To(BeNil())
			})
		})

		When("given module version does not exist", func() {
			It("returns no error", func() {
				err := repo.DeleteModuleVersion("com.example", "product", "go", "unknown")
				Expect(err).To(BeNil())
			})
		})

		When("given module version does exist", func() {
			It("returns no error", func() {
				err := repo.DeleteModuleVersion("com.example", "product", "go", "v1.0.0")
				Expect(err).To(BeNil())
			})
		})
	})

	Context("get module", func() {

		type args struct {
			namespace string
			name      string
			type_     string
			version   string
		}

		var (
			module *spec.Module
		)

		BeforeEach(func() {
			module = &spec.Module{
				Namespace: "com.example",
				Name:      "product",
				Type:      "go",
				Version: &spec.ModuleVersion{
					Name: "v1.0.0",
				},
			}

			Expect(repo.AddModule(module)).To(BeNil())
		})

		for _, tt := range []struct {
			name string
			args args
		}{
			{name: "namespace not known", args: args{namespace: "unknown", name: "product", type_: "go", version: "v1.0.0"}},
			{name: "name not known", args: args{namespace: "com.example", name: "unknown", type_: "go", version: "v1.0.0"}},
			{name: "type not known", args: args{namespace: "com.example", name: "product", type_: "unknown", version: "v1.0.0"}},
			{name: "version not known", args: args{namespace: "com.example", name: "product", type_: "go", version: "unknown"}},
		} {
			When(tt.name, func() {
				It("returns not found error", func() {
					m, err := repo.GetModule(tt.args.namespace, tt.args.name, tt.args.type_, tt.args.version)
					Expect(m).To(BeNil())
					Expect(err).To(MatchError("not found"))
				})
			})
		}

		When("module exists", func() {
			It("returns module and no error", func() {
				m, err := repo.GetModule("com.example", "product", "go", "v1.0.0")
				Expect(err).To(BeNil())
				Expect(proto.Equal(m, module)).To(BeTrue())
			})
		})
	})

	Context("list module namespaces", func() {

		When("no modules added", func() {
			It("returns empty namespace slice and no error", func() {
				namespaces, err := repo.ListModuleNamespaces()
				Expect(err).To(BeNil())
				Expect(namespaces).To(BeEmpty())
			})
		})

		When("modules added", func() {
			BeforeEach(func() {
				Expect(repo.AddModule(&spec.Module{
					Namespace: "com.example",
					Name:      "product",
					Type:      "go",
					Version: &spec.ModuleVersion{
						Name: "v1.0.0",
					},
				})).To(BeNil())
				Expect(repo.AddModule(&spec.Module{
					Namespace: "com.other",
					Name:      "customer",
					Type:      "go",
					Version: &spec.ModuleVersion{
						Name: "v2.0.0",
					},
				})).To(BeNil())
			})

			It("returns namespace slice and no error", func() {
				namespaces, err := repo.ListModuleNamespaces()
				Expect(err).To(BeNil())
				Expect(namespaces).To(HaveLen(2))
				Expect(namespaces).To(ContainElements("com.example", "com.other"))
			})
		})

	})

	Context("list module names", func() {

		When("no modules added", func() {
			It("returns empty name slice and no error", func() {
				names, err := repo.ListModuleNames("com.example")
				Expect(err).To(BeNil())
				Expect(names).To(BeEmpty())
			})
		})

		When("modules added", func() {
			BeforeEach(func() {
				Expect(repo.AddModule(&spec.Module{
					Namespace: "com.example",
					Name:      "product",
					Type:      "go",
					Version: &spec.ModuleVersion{
						Name: "v1.0.0",
					},
				})).To(BeNil())
				Expect(repo.AddModule(&spec.Module{
					Namespace: "com.example",
					Name:      "customer",
					Type:      "go",
					Version: &spec.ModuleVersion{
						Name: "v2.0.0",
					},
				})).To(BeNil())
			})

			It("returns name slice and no error", func() {
				namespaces, err := repo.ListModuleNames("com.example")
				Expect(err).To(BeNil())
				Expect(namespaces).To(HaveLen(2))
				Expect(namespaces).To(ContainElements("product", "customer"))
			})
		})

	})

	Context("list module types", func() {

		When("no modules added", func() {
			It("returns empty type slice and no error", func() {
				types, err := repo.ListModuleTypes("com.example", "product")
				Expect(err).To(BeNil())
				Expect(types).To(BeEmpty())
			})
		})

		When("modules added", func() {
			BeforeEach(func() {
				Expect(repo.AddModule(&spec.Module{
					Namespace: "com.example",
					Name:      "product",
					Type:      "go",
					Version: &spec.ModuleVersion{
						Name: "v1.0.0",
					},
				})).To(BeNil())
				Expect(repo.AddModule(&spec.Module{
					Namespace: "com.example",
					Name:      "product",
					Type:      "helm",
					Version: &spec.ModuleVersion{
						Name: "v1.0.0",
					},
				})).To(BeNil())
			})

			It("returns type slice and no error", func() {
				types, err := repo.ListModuleTypes("com.example", "product")
				Expect(err).To(BeNil())
				Expect(types).To(HaveLen(2))
				Expect(types).To(ContainElements("go", "helm"))
			})
		})

	})

	Context("list module versions", func() {

		When("no modules added", func() {
			It("returns empty version slice and no error", func() {
				versions, err := repo.ListModuleVersions("com.example", "product", "go")
				Expect(err).To(BeNil())
				Expect(versions).To(BeEmpty())
			})
		})

		When("modules added", func() {
			BeforeEach(func() {
				Expect(repo.AddModule(&spec.Module{
					Namespace: "com.example",
					Name:      "product",
					Type:      "go",
					Version: &spec.ModuleVersion{
						Name: "v1.0.0",
					},
				})).To(BeNil())
				Expect(repo.AddModule(&spec.Module{
					Namespace: "com.example",
					Name:      "product",
					Type:      "go",
					Version: &spec.ModuleVersion{
						Name: "v2.0.0",
					},
				})).To(BeNil())
			})

			It("returns version slice and no error", func() {
				versions, err := repo.ListModuleVersions("com.example", "product", "go")
				Expect(err).To(BeNil())
				Expect(versions).To(HaveLen(2))
				Expect(versions).To(ContainElements("v1.0.0", "v2.0.0"))
			})
		})

	})

})
