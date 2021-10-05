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

package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	specv1 "github.com/opendependency/go-spec/pkg/spec/v1"
	"google.golang.org/protobuf/proto"
)

var _ = Describe("unmarshal module from reader", func() {
	var (
		newReader func(string) io.Reader
		content   *strings.Builder

		module *specv1.Module
		err    error
	)

	BeforeEach(func() {
		newReader = func(s string) io.Reader {
			return strings.NewReader(s)
		}
		content = &strings.Builder{}

		module = &specv1.Module{}
	})

	JustBeforeEach(func() {
		err = unmarshalModuleFromReader(module, newReader(content.String()))
	})

	When("reader could not read all", func() {
		BeforeEach(func() {
			newReader = func(s string) io.Reader {
				return readerWithError{}
			}
		})

		It("should return error format not supported", func() {
			Expect(err).To(MatchError("could not read all: something"))
			Expect(proto.Equal(module, &specv1.Module{})).To(BeTrue())
		})
	})

	When("reader is empty", func() {
		BeforeEach(func() {
			content.Reset()
		})

		It("should return error format not supported", func() {
			Expect(err).To(MatchError("format not supported"))
			Expect(module).To(Equal(&specv1.Module{}))
		})
	})

	When("reader contains a json", func() {
		BeforeEach(func() {
			_, _ = content.WriteString(`{"namespace":"com.example.shop","name":"products","type":"go"}`)
		})

		It("should unmarshal module", func() {
			Expect(err).To(BeNil())

			Expect(proto.Equal(module, &specv1.Module{
				Namespace: "com.example.shop",
				Name:      "products",
				Type:      "go",
			})).To(BeTrue())
		})
	})

	When("reader contains a invalid json", func() {
		BeforeEach(func() {
			_, _ = content.WriteString(`invalid`)
		})

		It("should return error", func() {
			Expect(err).To(Not(BeNil()))
			Expect(err.Error()).To(ContainSubstring("could not unmarshal json:"))
			Expect(err.Error()).To(ContainSubstring("syntax error (line 1:1): unexpected token \"invalid\""))
			Expect(proto.Equal(module, &specv1.Module{})).To(BeTrue())
		})
	})

	When("reader contains a yaml", func() {
		BeforeEach(func() {
			_, _ = content.WriteString(`---
namespace: com.example.shop
name: products
type: go`)
		})

		It("should unmarshal module", func() {
			Expect(err).To(BeNil())

			Expect(proto.Equal(module, &specv1.Module{
				Namespace: "com.example.shop",
				Name:      "products",
				Type:      "go",
			})).To(BeTrue())
		})
	})
})

var _ = Describe("unmarshal module from file", func() {
	var (
		newFile            func() (string, error)
		writeContentToFile func(string, string) error
		content            *strings.Builder

		module *specv1.Module
		err    error
	)

	BeforeEach(func() {
		newFile = func() (string, error) {
			f, err := os.CreateTemp("", "module*")
			if err != nil {
				return "", fmt.Errorf("could not create temporary module file: %w", err)
			}
			defer func() {
				_ = f.Close()
			}()

			return f.Name(), nil
		}
		writeContentToFile = func(path string, content string) error {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, os.ModePerm)
			if err != nil {
				return fmt.Errorf("could not open file: %w", err)
			}
			defer func() {
				_ = f.Close()
			}()

			if _, err := f.WriteString(content); err != nil {
				return fmt.Errorf("could not write string to file: %w", err)
			}

			return nil
		}
		content = &strings.Builder{}

		module = &specv1.Module{}
	})

	JustBeforeEach(func() {
		f, ferr := newFile()
		if ferr != nil {
			Fail(ferr.Error())
		}
		if werr := writeContentToFile(f, content.String()); werr != nil {
			Fail(werr.Error())
		}

		err = unmarshalModuleFromFile(module, f)
	})

	When("file does not exist", func() {
		BeforeEach(func() {
			newFile = func() (string, error) {
				return "not-existing", nil
			}
			writeContentToFile = func(path string, content string) error {
				return nil
			}
		})

		It("should return error", func() {
			Expect(err).To(MatchError("file does not exist"))
			Expect(proto.Equal(module, &specv1.Module{})).To(BeTrue())
		})
	})

	When("reader is empty", func() {
		BeforeEach(func() {
			content.Reset()
		})

		It("should return error format not supported", func() {
			Expect(err).To(MatchError("format not supported"))
			Expect(module).To(Equal(&specv1.Module{}))
		})
	})

	When("reader contains a json", func() {
		BeforeEach(func() {
			_, _ = content.WriteString(`{"namespace":"com.example.shop","name":"products","type":"go"}`)
		})

		It("should unmarshal module", func() {
			Expect(err).To(BeNil())

			Expect(proto.Equal(module, &specv1.Module{
				Namespace: "com.example.shop",
				Name:      "products",
				Type:      "go",
			})).To(BeTrue())
		})
	})

	When("reader contains a invalid json", func() {
		BeforeEach(func() {
			_, _ = content.WriteString(`invalid`)
		})

		It("should return error", func() {
			Expect(err).To(Not(BeNil()))
			Expect(err.Error()).To(ContainSubstring("could not unmarshal json:"))
			Expect(err.Error()).To(ContainSubstring("syntax error (line 1:1): unexpected token \"invalid\""))
			Expect(proto.Equal(module, &specv1.Module{})).To(BeTrue())
		})
	})

	When("reader contains a yaml", func() {
		BeforeEach(func() {
			_, _ = content.WriteString(`---
namespace: com.example.shop
name: products
type: go`)
		})

		It("should unmarshal module", func() {
			Expect(err).To(BeNil())

			Expect(proto.Equal(module, &specv1.Module{
				Namespace: "com.example.shop",
				Name:      "products",
				Type:      "go",
			})).To(BeTrue())
		})
	})
})

var _ = Describe("parse module dependency", func() {
	var (
		dependency string

		parsedModuleDependency *specv1.ModuleDependency
		err                    error
	)

	BeforeEach(func() {
		dependency = ""
	})

	JustBeforeEach(func() {
		parsedModuleDependency, err = parseModuleDependency(dependency)
	})

	When("dependency is empty", func() {
		BeforeEach(func() {
			dependency = ""
		})

		It("should return error", func() {
			Expect(err).To(MatchError("must be of notation '<NAMESPACE>:<NAME>:<TYPE>:<VERSION>'"))
			Expect(parsedModuleDependency).To(BeNil())
		})
	})

	When("dependency does not match notation", func() {
		BeforeEach(func() {
			dependency = "a:b"
		})

		It("should return error", func() {
			Expect(err).To(MatchError("must be of notation '<NAMESPACE>:<NAME>:<TYPE>:<VERSION>'"))
			Expect(parsedModuleDependency).To(BeNil())
		})
	})

	When("dependency match notation", func() {
		BeforeEach(func() {
			dependency = "a:b:c:d"
		})

		It("should return parsed module dependency", func() {
			Expect(err).To(BeNil())
			Expect(proto.Equal(parsedModuleDependency, &specv1.ModuleDependency{
				Namespace: "a",
				Name:      "b",
				Type:      "c",
				Version:   "d",
			})).To(BeTrue())
		})
	})
})

type readerWithError struct{}

func (r readerWithError) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("something")
}
