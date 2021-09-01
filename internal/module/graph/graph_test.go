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

package graph

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	spec "github.com/opendependency/go-spec/pkg/spec/v1"
)

var _ = Describe("graph", func() {

	var (
		m AdjacentMatrix
		g *graph
	)

	BeforeEach(func() {
		m = NewInMemoryAdjacentMatrix()
		g = NewGraph(m)
	})

	Context("add module", func() {

		When("module is nil", func() {
			It("returns an error", func() {
				err := g.AddModule(nil)

				Expect(err).To(MatchError("module must not be nil"))
			})
		})

		When("module is invalid", func() {
			It("returns an error", func() {
				err := g.AddModule(&spec.Module{})

				Expect(err).To(MatchError("module validation failed: namespace: must have at least 1 characters"))
			})
		})

		When("module has no dependencies", func() {
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
					Dependencies: nil,
				}
			})

			It("returns no error", func() {
				err := g.AddModule(module)

				Expect(err).To(BeNil())
			})

			It("adds no edges to adjacent matrix", func() {
				_ = g.AddModule(module)

				Expect(m.NumberOfEdges(dependsOnEdge)).To(Equal(0))
				Expect(m.NumberOfEdges(usedByEdge)).To(Equal(0))
				Expect(m.NumberOfEdges(requiredForEdge)).To(Equal(0))
				Expect(m.NumberOfEdges(requireEdge)).To(Equal(0))
			})
		})

		When("module has one upstream dependency", func() {
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
					Dependencies: []*spec.ModuleDependency{
						{
							Namespace: "com.example",
							Name:      "lib",
							Type:      "go",
							Version:   "v1.2.3",
						},
					},
				}
			})

			It("returns no error", func() {
				err := g.AddModule(module)

				Expect(err).To(BeNil())
			})

			It("adds a depend-on edge from product to lib", func() {
				_ = g.AddModule(module)

				Expect(m.NumberOfEdges(dependsOnEdge)).To(Equal(1))
				Expect(m.Get(dependsOnEdge, Vertex{
					Namespace: "com.example",
					Name:      "product",
					Type:      "go",
					Version:   "v1.0.0",
				})).To(ContainElement(Vertex{
					Namespace: "com.example",
					Name:      "lib",
					Type:      "go",
					Version:   "v1.2.3",
				}))
			})

			It("adds an used-by edge from product to lib", func() {
				_ = g.AddModule(module)

				Expect(m.NumberOfEdges(usedByEdge)).To(Equal(1))
				Expect(m.Get(usedByEdge, Vertex{
					Namespace: "com.example",
					Name:      "lib",
					Type:      "go",
					Version:   "v1.2.3",
				})).To(ContainElement(Vertex{
					Namespace: "com.example",
					Name:      "product",
					Type:      "go",
					Version:   "v1.0.0",
				}))
			})

			It("does not add a required-for edge", func() {
				_ = g.AddModule(module)

				Expect(m.NumberOfEdges(requiredForEdge)).To(Equal(0))
			})

			It("does not add a require edge", func() {
				_ = g.AddModule(module)

				Expect(m.NumberOfEdges(requireEdge)).To(Equal(0))
			})
		})

		When("module has one downstream dependency", func() {
			var (
				module *spec.Module
			)

			BeforeEach(func() {
				downstreamDirection := spec.DependencyDirection_DOWNSTREAM
				module = &spec.Module{
					Namespace: "com.example",
					Name:      "product",
					Type:      "go",
					Version: &spec.ModuleVersion{
						Name: "v1.0.0",
					},
					Dependencies: []*spec.ModuleDependency{
						{
							Namespace: "com.example",
							Name:      "product",
							Type:      "protobuf",
							Version:   "v1.8.9",
							Direction: &downstreamDirection,
						},
					},
				}
			})

			It("returns no error", func() {
				err := g.AddModule(module)

				Expect(err).To(BeNil())
			})

			It("adds a required-for edge from product go to product protobuf", func() {
				_ = g.AddModule(module)

				Expect(m.NumberOfEdges(requiredForEdge)).To(Equal(1))
				Expect(m.Get(requiredForEdge, Vertex{
					Namespace: "com.example",
					Name:      "product",
					Type:      "go",
					Version:   "v1.0.0",
				})).To(ContainElement(Vertex{
					Namespace: "com.example",
					Name:      "product",
					Type:      "protobuf",
					Version:   "v1.8.9",
				}))
			})

			It("adds a require edge from product protobuf to product go", func() {
				_ = g.AddModule(module)

				Expect(m.NumberOfEdges(requireEdge)).To(Equal(1))
				Expect(m.Get(requireEdge, Vertex{
					Namespace: "com.example",
					Name:      "product",
					Type:      "protobuf",
					Version:   "v1.8.9",
				})).To(ContainElement(Vertex{
					Namespace: "com.example",
					Name:      "product",
					Type:      "go",
					Version:   "v1.0.0",
				}))
			})

			It("does not add a depend-on edge", func() {
				_ = g.AddModule(module)

				Expect(m.NumberOfEdges(dependsOnEdge)).To(Equal(0))
			})

			It("does not add an used-by edge", func() {
				_ = g.AddModule(module)

				Expect(m.NumberOfEdges(usedByEdge)).To(Equal(0))
			})
		})

	})

	Context("traverse breadth first search", func() {
		var (
			startVertex Vertex
		)

		BeforeEach(func() {
			startVertex = Vertex{
				Namespace: "com.example",
				Name:      "product",
				Type:      "go",
				Version:   "v1.0.0",
			}
		})

		When("adjacent matrix is empty", func() {
			It("return start vertex as parent", func() {
				called := false
				g.traverseBFS("my-edge", startVertex, func(p Vertex, v []Vertex) bool {
					called = true
					Expect(p).To(Equal(startVertex))
					return false
				})
				Expect(called).To(BeTrue())
			})

			It("return an empty vertex slice as children", func() {
				called := false
				g.traverseBFS("my-edge", startVertex, func(p Vertex, v []Vertex) bool {
					called = true
					Expect(v).To(BeEmpty())
					return false
				})
				Expect(called).To(BeTrue())
			})

			It("is only called once", func() {
				called := 0
				g.traverseBFS("my-edge", startVertex, func(p Vertex, v []Vertex) bool {
					called++
					return true
				})
				Expect(called).To(Equal(1))
			})
		})

		When("adjacent matrix is not empty", func() {

			type fnArgs struct {
				p Vertex
				v []Vertex
			}

			var (
				expectedFnCalls []fnArgs
			)

			BeforeEach(func() {
				timeLibGo := Vertex{
					Namespace: "com.example",
					Name:      "time-lib",
					Type:      "go",
					Version:   "v3.1.0",
				}
				utilLibGo := Vertex{
					Namespace: "com.example",
					Name:      "util-lib",
					Type:      "go",
					Version:   "v5.0.0",
				}
				pricingProtobuf := Vertex{
					Namespace: "com.example",
					Name:      "pricing",
					Type:      "protobuf",
					Version:   "v3.0.0",
				}

				m.AddEdges("my-edge", startVertex, []Vertex{
					utilLibGo,
					pricingProtobuf,
				})

				m.AddEdges("my-edge", utilLibGo, []Vertex{
					timeLibGo,
				})

				expectedFnCalls = []fnArgs{
					{startVertex, []Vertex{utilLibGo, pricingProtobuf}},
					{utilLibGo, []Vertex{timeLibGo}},
					{pricingProtobuf, []Vertex{}},
					{timeLibGo, []Vertex{}},
				}
			})

			It("call the function with start vertex as parent", func() {
				called := false
				g.traverseBFS("my-edge", startVertex, func(p Vertex, v []Vertex) bool {
					called = true
					Expect(p).To(Equal(startVertex))
					return false
				})
				Expect(called).To(BeTrue())
			})

			It("call the function as expected", func() {
				called := 0
				g.traverseBFS("my-edge", startVertex, func(p Vertex, v []Vertex) bool {
					if called >= len(expectedFnCalls) {
						Fail("called too much")
					}

					args := expectedFnCalls[called]
					Expect(p).To(Equal(args.p))
					Expect(v).To(ContainElements(args.v))
					called++
					return true
				})
				Expect(called).To(Equal(len(expectedFnCalls)))
			})
		})
	})

	Context("traverse depth first search", func() {
		var (
			startVertex Vertex
		)

		BeforeEach(func() {
			startVertex = Vertex{
				Namespace: "com.example",
				Name:      "product",
				Type:      "go",
				Version:   "v1.0.0",
			}
		})

		When("adjacent matrix is empty", func() {
			It("does call function", func() {
				called := false
				g.traverseDFS("my-edge", startVertex, func(p Vertex, v Vertex) bool {
					Expect(p).To(Equal(Vertex{}))
					Expect(v).To(Equal(startVertex))
					called = true
					return false
				})
				Expect(called).To(BeTrue())
			})
		})

		When("adjacent matrix is not empty", func() {

			type fnArgs struct {
				p Vertex
				v Vertex
			}

			var (
				expectedFnCalls []fnArgs
			)

			BeforeEach(func() {
				timeLibGo := Vertex{
					Namespace: "com.example",
					Name:      "time-lib",
					Type:      "go",
					Version:   "v3.1.0",
				}
				utilLibGo := Vertex{
					Namespace: "com.example",
					Name:      "util-lib",
					Type:      "go",
					Version:   "v5.0.0",
				}
				pricingProtobuf := Vertex{
					Namespace: "com.example",
					Name:      "pricing",
					Type:      "protobuf",
					Version:   "v3.0.0",
				}

				m.AddEdges("my-edge", startVertex, []Vertex{
					utilLibGo,
					pricingProtobuf,
				})

				m.AddEdges("my-edge", utilLibGo, []Vertex{
					timeLibGo,
				})

				expectedFnCalls = []fnArgs{
					{Vertex{}, startVertex},
					{startVertex, pricingProtobuf},
					{startVertex, utilLibGo},
					{utilLibGo, timeLibGo},
				}
			})

			It("call the function with empty vertex as parent", func() {
				called := false
				g.traverseDFS("my-edge", startVertex, func(p Vertex, v Vertex) bool {
					called = true
					Expect(p).To(Equal(Vertex{}))
					Expect(v).To(Equal(startVertex))
					return false
				})
				Expect(called).To(BeTrue())
			})

			It("call the function as expected", func() {
				called := 0
				g.traverseDFS("my-edge", startVertex, func(p Vertex, v Vertex) bool {
					if called >= len(expectedFnCalls) {
						Fail("called too much")
					}

					args := expectedFnCalls[called]
					Expect(p).To(Equal(args.p))
					Expect(v).To(Equal(args.v))
					called++
					return true
				})
				Expect(called).To(Equal(len(expectedFnCalls)))
			})
		})
	})

	Context("traverse * edges *", func() {
		BeforeEach(func() {
			// The desired graph looks like the following:
			//
			//       (com.example:product:helm:v1.5.0)                            (com.example:order:helm:v2.3.8)
			//             |               ^                                              |               ^
			//             |               |                                              |               |
			//          depend-on        used-by                                       depend-on        used-by
			//             |               |                                              |               |
			//             v               |                                              v               |
			//  (com.example:product:container-image:v1.5.0)               (com.example:order:container-image:v2.3.8)
			//             |               ^                                              |               ^
			//             |               |                                              |               |
			//         depend-on        used-by                                        depend-on        used-by
			//             |               |                                              |               |
			//             v               |                                              v               |
			//     (com.example:product:go:v1.5.0)                                 (com.example:order:go:v2.3.8)----depend-on---->(com.example:utils:go:v4.3.1)
			//             ^                 \                                            /               ^   ^                    /
			//              \                 \                                          /               /     \                  /
			//               \            required-for                              depend-on        used-by    -----used-by------
			//             require              \                                      /               /
			//                 \                 \                                    /               /
			//                  \                 v                                  v               /
			//                   -----------------(com.example:product:protobuf:v1.0.0)--------------
			//
			downstreamDirection := spec.DependencyDirection_DOWNSTREAM

			for _, mod := range []*spec.Module{
				{
					Namespace: "com.example",
					Name:      "product",
					Type:      "helm",
					Version:   &spec.ModuleVersion{Name: "v1.5.0"},
					Dependencies: []*spec.ModuleDependency{
						{
							Namespace: "com.example",
							Name:      "product",
							Type:      "container-image",
							Version:   "v1.5.0",
						},
					},
				},
				{
					Namespace: "com.example",
					Name:      "product",
					Type:      "container-image",
					Version:   &spec.ModuleVersion{Name: "v1.5.0"},
					Dependencies: []*spec.ModuleDependency{
						{
							Namespace: "com.example",
							Name:      "product",
							Type:      "go",
							Version:   "v1.5.0",
						},
					},
				},
				{
					Namespace: "com.example",
					Name:      "product",
					Type:      "go",
					Version:   &spec.ModuleVersion{Name: "v1.5.0"},
					Dependencies: []*spec.ModuleDependency{
						{
							Namespace: "com.example",
							Name:      "product",
							Type:      "protobuf",
							Version:   "v1.0.0",
							Direction: &downstreamDirection,
						},
					},
				},
				{
					Namespace: "com.example",
					Name:      "product",
					Type:      "protobuf",
					Version:   &spec.ModuleVersion{Name: "v1.0.0"},
				},

				{
					Namespace: "com.example",
					Name:      "order",
					Type:      "helm",
					Version:   &spec.ModuleVersion{Name: "v2.3.8"},
					Dependencies: []*spec.ModuleDependency{
						{
							Namespace: "com.example",
							Name:      "order",
							Type:      "container-image",
							Version:   "v2.3.8",
						},
					},
				},
				{
					Namespace: "com.example",
					Name:      "order",
					Type:      "container-image",
					Version:   &spec.ModuleVersion{Name: "v2.3.8"},
					Dependencies: []*spec.ModuleDependency{
						{
							Namespace: "com.example",
							Name:      "order",
							Type:      "go",
							Version:   "v2.3.8",
						},
					},
				},
				{
					Namespace: "com.example",
					Name:      "order",
					Type:      "go",
					Version:   &spec.ModuleVersion{Name: "v2.3.8"},
					Dependencies: []*spec.ModuleDependency{
						{
							Namespace: "com.example",
							Name:      "product",
							Type:      "protobuf",
							Version:   "v1.0.0",
						},
						{
							Namespace: "com.example",
							Name:      "utils",
							Type:      "go",
							Version:   "v4.3.1",
						},
					},
				},
			} {
				if err := g.AddModule(mod); err != nil {
					Fail(err.Error())
				}
			}
		})

		Context("traverse depends-on edges bfs", func() {
			type fnArgs struct {
				p Vertex
				v []Vertex
			}

			var (
				startVertex     Vertex
				expectedFnCalls []fnArgs
			)

			When("start vertex is com.example:product:helm:v1.5.0", func() {

				BeforeEach(func() {
					startVertex = Vertex{Namespace: "com.example", Name: "product", Type: "helm", Version: "v1.5.0"}
					expectedFnCalls = []fnArgs{
						{p: Vertex{Namespace: "com.example", Name: "product", Type: "helm", Version: "v1.5.0"}, v: []Vertex{{Namespace: "com.example", Name: "product", Type: "container-image", Version: "v1.5.0"}}},
						{p: Vertex{Namespace: "com.example", Name: "product", Type: "container-image", Version: "v1.5.0"}, v: []Vertex{{Namespace: "com.example", Name: "product", Type: "go", Version: "v1.5.0"}}},
						{p: Vertex{Namespace: "com.example", Name: "product", Type: "go", Version: "v1.5.0"}, v: []Vertex{}},
					}
				})

				It("call the function as expected", func() {
					called := 0
					g.TraverseDependOnEdgesBFS(startVertex, func(p Vertex, v []Vertex) bool {
						if called >= len(expectedFnCalls) {
							Fail("called too much")
						}

						args := expectedFnCalls[called]
						Expect(p).To(Equal(args.p))
						Expect(v).To(HaveLen(len(args.v)))
						Expect(v).To(ContainElements(args.v))
						called++
						return true
					})
					Expect(called).To(Equal(len(expectedFnCalls)))
				})
			})

			When("start vertex is com.example:order:helm:v2.3.8", func() {

				BeforeEach(func() {
					startVertex = Vertex{Namespace: "com.example", Name: "order", Type: "helm", Version: "v2.3.8"}
					expectedFnCalls = []fnArgs{
						{p: Vertex{Namespace: "com.example", Name: "order", Type: "helm", Version: "v2.3.8"}, v: []Vertex{{Namespace: "com.example", Name: "order", Type: "container-image", Version: "v2.3.8"}}},
						{p: Vertex{Namespace: "com.example", Name: "order", Type: "container-image", Version: "v2.3.8"}, v: []Vertex{{Namespace: "com.example", Name: "order", Type: "go", Version: "v2.3.8"}}},
						{p: Vertex{Namespace: "com.example", Name: "order", Type: "go", Version: "v2.3.8"}, v: []Vertex{{Namespace: "com.example", Name: "product", Type: "protobuf", Version: "v1.0.0"}, {Namespace: "com.example", Name: "utils", Type: "go", Version: "v4.3.1"}}},
						{p: Vertex{Namespace: "com.example", Name: "product", Type: "protobuf", Version: "v1.0.0"}, v: []Vertex{}},
						{p: Vertex{Namespace: "com.example", Name: "utils", Type: "go", Version: "v4.3.1"}, v: []Vertex{}},
					}
				})

				It("call the function as expected", func() {
					called := 0
					g.TraverseDependOnEdgesBFS(startVertex, func(p Vertex, v []Vertex) bool {
						if called >= len(expectedFnCalls) {
							Fail("called too much")
						}

						args := expectedFnCalls[called]
						Expect(p).To(Equal(args.p))
						Expect(v).To(HaveLen(len(args.v)))
						Expect(v).To(ContainElements(args.v))
						called++
						return true
					})
					Expect(called).To(Equal(len(expectedFnCalls)))
				})
			})

		})

		Context("traverse depends-on edges dfs", func() {
			type fnArgs struct {
				p Vertex
				v Vertex
			}

			var (
				startVertex     Vertex
				expectedFnCalls []fnArgs
			)

			When("start vertex is com.example:product:helm:v1.5.0", func() {

				BeforeEach(func() {
					startVertex = Vertex{Namespace: "com.example", Name: "product", Type: "helm", Version: "v1.5.0"}
					expectedFnCalls = []fnArgs{
						{p: Vertex{}, v: Vertex{Namespace: "com.example", Name: "product", Type: "helm", Version: "v1.5.0"}},
						{p: Vertex{Namespace: "com.example", Name: "product", Type: "helm", Version: "v1.5.0"}, v: Vertex{Namespace: "com.example", Name: "product", Type: "container-image", Version: "v1.5.0"}},
						{p: Vertex{Namespace: "com.example", Name: "product", Type: "container-image", Version: "v1.5.0"}, v: Vertex{Namespace: "com.example", Name: "product", Type: "go", Version: "v1.5.0"}},
					}
				})

				It("call the function as expected", func() {
					called := 0
					g.TraverseDependOnEdgesDFS(startVertex, func(p Vertex, v Vertex) bool {
						if called >= len(expectedFnCalls) {
							Fail("called too much")
						}

						args := expectedFnCalls[called]
						Expect(p).To(Equal(args.p))
						Expect(v).To(Equal(args.v))
						called++
						return true
					})
					Expect(called).To(Equal(len(expectedFnCalls)))
				})
			})

			When("start vertex is com.example:order:helm:v2.3.8", func() {

				BeforeEach(func() {
					startVertex = Vertex{Namespace: "com.example", Name: "order", Type: "helm", Version: "v2.3.8"}
					expectedFnCalls = []fnArgs{
						{p: Vertex{}, v: Vertex{Namespace: "com.example", Name: "order", Type: "helm", Version: "v2.3.8"}},
						{p: Vertex{Namespace: "com.example", Name: "order", Type: "helm", Version: "v2.3.8"}, v: Vertex{Namespace: "com.example", Name: "order", Type: "container-image", Version: "v2.3.8"}},
						{p: Vertex{Namespace: "com.example", Name: "order", Type: "container-image", Version: "v2.3.8"}, v: Vertex{Namespace: "com.example", Name: "order", Type: "go", Version: "v2.3.8"}},
						{p: Vertex{Namespace: "com.example", Name: "order", Type: "go", Version: "v2.3.8"}, v: Vertex{Namespace: "com.example", Name: "utils", Type: "go", Version: "v4.3.1"}},
						{p: Vertex{Namespace: "com.example", Name: "order", Type: "go", Version: "v2.3.8"}, v: Vertex{Namespace: "com.example", Name: "product", Type: "protobuf", Version: "v1.0.0"}},
					}
				})

				It("call the function as expected", func() {
					called := 0
					g.TraverseDependOnEdgesDFS(startVertex, func(p Vertex, v Vertex) bool {
						if called >= len(expectedFnCalls) {
							Fail("called too much")
						}

						args := expectedFnCalls[called]
						Expect(p).To(Equal(args.p))
						Expect(v).To(Equal(args.v))
						called++
						return true
					})
					Expect(called).To(Equal(len(expectedFnCalls)))
				})
			})

		})

		Context("traverse used-by edges bfs", func() {
			type fnArgs struct {
				p Vertex
				v []Vertex
			}

			var (
				startVertex     Vertex
				expectedFnCalls []fnArgs
			)

			When("start vertex is com.example:product:go:v1.5.0", func() {

				BeforeEach(func() {
					startVertex = Vertex{Namespace: "com.example", Name: "product", Type: "go", Version: "v1.5.0"}
					expectedFnCalls = []fnArgs{
						{p: Vertex{Namespace: "com.example", Name: "product", Type: "go", Version: "v1.5.0"}, v: []Vertex{{Namespace: "com.example", Name: "product", Type: "container-image", Version: "v1.5.0"}}},
						{p: Vertex{Namespace: "com.example", Name: "product", Type: "container-image", Version: "v1.5.0"}, v: []Vertex{{Namespace: "com.example", Name: "product", Type: "helm", Version: "v1.5.0"}}},
						{p: Vertex{Namespace: "com.example", Name: "product", Type: "helm", Version: "v1.5.0"}, v: []Vertex{}},
					}
				})

				It("call the function as expected", func() {
					called := 0
					g.TraverseUsedByEdgesBFS(startVertex, func(p Vertex, v []Vertex) bool {
						if called >= len(expectedFnCalls) {
							Fail("called too much")
						}

						args := expectedFnCalls[called]
						Expect(p).To(Equal(args.p))
						Expect(v).To(HaveLen(len(args.v)))
						Expect(v).To(ContainElements(args.v))
						called++
						return true
					})
					Expect(called).To(Equal(len(expectedFnCalls)))
				})
			})

			When("start vertex is com.example:product:protobuf:v1.0.0", func() {

				BeforeEach(func() {
					startVertex = Vertex{Namespace: "com.example", Name: "product", Type: "protobuf", Version: "v1.0.0"}
					expectedFnCalls = []fnArgs{
						{p: Vertex{Namespace: "com.example", Name: "product", Type: "protobuf", Version: "v1.0.0"}, v: []Vertex{{Namespace: "com.example", Name: "order", Type: "go", Version: "v2.3.8"}}},
						{p: Vertex{Namespace: "com.example", Name: "order", Type: "go", Version: "v2.3.8"}, v: []Vertex{{Namespace: "com.example", Name: "order", Type: "container-image", Version: "v2.3.8"}}},
						{p: Vertex{Namespace: "com.example", Name: "order", Type: "container-image", Version: "v2.3.8"}, v: []Vertex{{Namespace: "com.example", Name: "order", Type: "helm", Version: "v2.3.8"}}},
						{p: Vertex{Namespace: "com.example", Name: "order", Type: "helm", Version: "v2.3.8"}, v: []Vertex{}},
					}
				})

				It("call the function as expected", func() {
					called := 0
					g.TraverseUsedByEdgesBFS(startVertex, func(p Vertex, v []Vertex) bool {
						if called >= len(expectedFnCalls) {
							Fail("called too much")
						}

						args := expectedFnCalls[called]
						Expect(p).To(Equal(args.p))
						Expect(v).To(HaveLen(len(args.v)))
						Expect(v).To(ContainElements(args.v))
						called++
						return true
					})
					Expect(called).To(Equal(len(expectedFnCalls)))
				})
			})

		})

		Context("traverse used-by edges dfs", func() {
			type fnArgs struct {
				p Vertex
				v Vertex
			}

			var (
				startVertex     Vertex
				expectedFnCalls []fnArgs
			)

			When("start vertex is com.example:product:go:v1.5.0", func() {

				BeforeEach(func() {
					startVertex = Vertex{Namespace: "com.example", Name: "product", Type: "go", Version: "v1.5.0"}
					expectedFnCalls = []fnArgs{
						{p: Vertex{}, v: Vertex{Namespace: "com.example", Name: "product", Type: "go", Version: "v1.5.0"}},
						{p: Vertex{Namespace: "com.example", Name: "product", Type: "go", Version: "v1.5.0"}, v: Vertex{Namespace: "com.example", Name: "product", Type: "container-image", Version: "v1.5.0"}},
						{p: Vertex{Namespace: "com.example", Name: "product", Type: "container-image", Version: "v1.5.0"}, v: Vertex{Namespace: "com.example", Name: "product", Type: "helm", Version: "v1.5.0"}},
					}
				})

				It("call the function as expected", func() {
					called := 0
					g.TraverseUsedByEdgesDFS(startVertex, func(p Vertex, v Vertex) bool {
						if called >= len(expectedFnCalls) {
							Fail("called too much")
						}

						args := expectedFnCalls[called]
						Expect(p).To(Equal(args.p))
						Expect(v).To(Equal(args.v))
						called++
						return true
					})
					Expect(called).To(Equal(len(expectedFnCalls)))
				})
			})

			When("start vertex is com.example:product:protobuf:v1.0.0", func() {

				BeforeEach(func() {
					startVertex = Vertex{Namespace: "com.example", Name: "product", Type: "protobuf", Version: "v1.0.0"}
					expectedFnCalls = []fnArgs{
						{p: Vertex{}, v: Vertex{Namespace: "com.example", Name: "product", Type: "protobuf", Version: "v1.0.0"}},
						{p: Vertex{Namespace: "com.example", Name: "product", Type: "protobuf", Version: "v1.0.0"}, v: Vertex{Namespace: "com.example", Name: "order", Type: "go", Version: "v2.3.8"}},
						{p: Vertex{Namespace: "com.example", Name: "order", Type: "go", Version: "v2.3.8"}, v: Vertex{Namespace: "com.example", Name: "order", Type: "container-image", Version: "v2.3.8"}},
						{p: Vertex{Namespace: "com.example", Name: "order", Type: "container-image", Version: "v2.3.8"}, v: Vertex{Namespace: "com.example", Name: "order", Type: "helm", Version: "v2.3.8"}},
					}
				})

				It("call the function as expected", func() {
					called := 0
					g.TraverseUsedByEdgesDFS(startVertex, func(p Vertex, v Vertex) bool {
						if called >= len(expectedFnCalls) {
							Fail("called too much")
						}

						args := expectedFnCalls[called]
						Expect(p).To(Equal(args.p))
						Expect(v).To(Equal(args.v))
						called++
						return true
					})
					Expect(called).To(Equal(len(expectedFnCalls)))
				})
			})

		})

		Context("traverse required-for edges bfs", func() {
			type fnArgs struct {
				p Vertex
				v []Vertex
			}

			var (
				startVertex     Vertex
				expectedFnCalls []fnArgs
			)

			When("start vertex is com.example:product:go:v1.5.0", func() {

				BeforeEach(func() {
					startVertex = Vertex{Namespace: "com.example", Name: "product", Type: "go", Version: "v1.5.0"}
					expectedFnCalls = []fnArgs{
						{p: Vertex{Namespace: "com.example", Name: "product", Type: "go", Version: "v1.5.0"}, v: []Vertex{{Namespace: "com.example", Name: "product", Type: "protobuf", Version: "v1.0.0"}}},
						{p: Vertex{Namespace: "com.example", Name: "product", Type: "protobuf", Version: "v1.0.0"}, v: []Vertex{}},
					}
				})

				It("call the function as expected", func() {
					called := 0
					g.TraverseRequiredForEdgesBFS(startVertex, func(p Vertex, v []Vertex) bool {
						if called >= len(expectedFnCalls) {
							Fail("called too much")
						}

						args := expectedFnCalls[called]
						Expect(p).To(Equal(args.p))
						Expect(v).To(HaveLen(len(args.v)))
						Expect(v).To(ContainElements(args.v))
						called++
						return true
					})
					Expect(called).To(Equal(len(expectedFnCalls)))
				})
			})

		})

		Context("traverse required-for edges dfs", func() {
			type fnArgs struct {
				p Vertex
				v Vertex
			}

			var (
				startVertex     Vertex
				expectedFnCalls []fnArgs
			)

			When("start vertex is com.example:product:go:v1.5.0", func() {

				BeforeEach(func() {
					startVertex = Vertex{Namespace: "com.example", Name: "product", Type: "go", Version: "v1.5.0"}
					expectedFnCalls = []fnArgs{
						{p: Vertex{}, v: Vertex{Namespace: "com.example", Name: "product", Type: "go", Version: "v1.5.0"}},
						{p: Vertex{Namespace: "com.example", Name: "product", Type: "go", Version: "v1.5.0"}, v: Vertex{Namespace: "com.example", Name: "product", Type: "protobuf", Version: "v1.0.0"}},
					}
				})

				It("call the function as expected", func() {
					called := 0
					g.TraverseRequiredForEdgesDFS(startVertex, func(p Vertex, v Vertex) bool {
						if called >= len(expectedFnCalls) {
							Fail("called too much")
						}

						args := expectedFnCalls[called]
						Expect(p).To(Equal(args.p))
						Expect(v).To(Equal(args.v))
						called++
						return true
					})
					Expect(called).To(Equal(len(expectedFnCalls)))
				})
			})

		})

		Context("traverse require edges bfs", func() {
			type fnArgs struct {
				p Vertex
				v []Vertex
			}

			var (
				startVertex     Vertex
				expectedFnCalls []fnArgs
			)

			When("start vertex is com.example:product:protobuf:v1.0.0", func() {

				BeforeEach(func() {
					startVertex = Vertex{Namespace: "com.example", Name: "product", Type: "protobuf", Version: "v1.0.0"}
					expectedFnCalls = []fnArgs{
						{p: Vertex{Namespace: "com.example", Name: "product", Type: "protobuf", Version: "v1.0.0"}, v: []Vertex{{Namespace: "com.example", Name: "product", Type: "go", Version: "v1.5.0"}}},
						{p: Vertex{Namespace: "com.example", Name: "product", Type: "go", Version: "v1.5.0"}, v: []Vertex{}},
					}
				})

				It("call the function as expected", func() {
					called := 0
					g.TraverseRequireEdgesBFS(startVertex, func(p Vertex, v []Vertex) bool {
						if called >= len(expectedFnCalls) {
							Fail("called too much")
						}

						args := expectedFnCalls[called]
						Expect(p).To(Equal(args.p))
						Expect(v).To(HaveLen(len(args.v)))
						Expect(v).To(ContainElements(args.v))
						called++
						return true
					})
					Expect(called).To(Equal(len(expectedFnCalls)))
				})
			})

		})

		Context("traverse require edges dfs", func() {
			type fnArgs struct {
				p Vertex
				v Vertex
			}

			var (
				startVertex     Vertex
				expectedFnCalls []fnArgs
			)

			When("start vertex is com.example:product:protobuf:v1.0.0", func() {

				BeforeEach(func() {
					startVertex = Vertex{Namespace: "com.example", Name: "product", Type: "protobuf", Version: "v1.0.0"}
					expectedFnCalls = []fnArgs{
						{p: Vertex{}, v: Vertex{Namespace: "com.example", Name: "product", Type: "protobuf", Version: "v1.0.0"}},
						{p: Vertex{Namespace: "com.example", Name: "product", Type: "protobuf", Version: "v1.0.0"}, v: Vertex{Namespace: "com.example", Name: "product", Type: "go", Version: "v1.5.0"}},
					}
				})

				It("call the function as expected", func() {
					called := 0
					g.TraverseRequireEdgesDFS(startVertex, func(p Vertex, v Vertex) bool {
						if called >= len(expectedFnCalls) {
							Fail("called too much")
						}

						args := expectedFnCalls[called]
						Expect(p).To(Equal(args.p))
						Expect(v).To(Equal(args.v))
						called++
						return true
					})
					Expect(called).To(Equal(len(expectedFnCalls)))
				})
			})

		})
	})
})
