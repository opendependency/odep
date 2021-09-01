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
)

var _ = Describe("in-memory adjacent matrix", func() {

	var (
		matrix *inMemoryAdjacentMatrix
	)

	BeforeEach(func() {
		matrix = NewInMemoryAdjacentMatrix()
	})

	Context("add edge", func() {
		When("name is empty", func() {
			It("adds an edge", func() {
				matrix.AddEdge("", Vertex{"a", "b", "c", "d"}, Vertex{"e", "f", "g", "h"})

				Expect(matrix.m).To(HaveLen(1))
				Expect(matrix.m[""]).To(HaveLen(1))
				Expect(matrix.m[""][Vertex{"a", "b", "c", "d"}]).To(HaveLen(1))
			})
		})

		When("name is not empty", func() {
			It("adds an edge", func() {
				matrix.AddEdge("upstream", Vertex{"a", "b", "c", "d"}, Vertex{"e", "f", "g", "h"})

				Expect(matrix.m).To(HaveLen(1))
				Expect(matrix.m["upstream"]).To(HaveLen(1))
				Expect(matrix.m["upstream"][Vertex{"a", "b", "c", "d"}]).To(HaveLen(1))
			})
		})

		When("parent vertex is empty", func() {
			It("adds an edge", func() {
				matrix.AddEdge("upstream", Vertex{}, Vertex{"e", "f", "g", "h"})

				Expect(matrix.m).To(HaveLen(1))
				Expect(matrix.m["upstream"]).To(HaveLen(1))
				Expect(matrix.m["upstream"][Vertex{}]).To(HaveLen(1))
			})
		})

		When("child vertex is empty", func() {
			It("adds an edge", func() {
				matrix.AddEdge("upstream", Vertex{"a", "b", "c", "d"}, Vertex{})

				Expect(matrix.m).To(HaveLen(1))
				Expect(matrix.m["upstream"]).To(HaveLen(1))
				Expect(matrix.m["upstream"][Vertex{"a", "b", "c", "d"}]).To(HaveLen(1))
			})
		})
	})

	Context("add edges", func() {
		When("name is empty", func() {
			It("adds an edge", func() {
				matrix.AddEdges("", Vertex{"a", "b", "c", "d"}, []Vertex{{"e", "f", "g", "h"}, {"i", "j", "k", "l"}})

				Expect(matrix.m).To(HaveLen(1))
				Expect(matrix.m[""]).To(HaveLen(1))
				Expect(matrix.m[""][Vertex{"a", "b", "c", "d"}]).To(HaveLen(2))
			})
		})

		When("name is not empty", func() {
			It("adds an edge", func() {
				matrix.AddEdges("upstream", Vertex{"a", "b", "c", "d"}, []Vertex{{"e", "f", "g", "h"}, {"i", "j", "k", "l"}})

				Expect(matrix.m).To(HaveLen(1))
				Expect(matrix.m["upstream"]).To(HaveLen(1))
				Expect(matrix.m["upstream"][Vertex{"a", "b", "c", "d"}]).To(HaveLen(2))
			})
		})

		When("parent vertex is empty", func() {
			It("adds an edge", func() {
				matrix.AddEdges("upstream", Vertex{}, []Vertex{{"e", "f", "g", "h"}, {"i", "j", "k", "l"}})

				Expect(matrix.m).To(HaveLen(1))
				Expect(matrix.m["upstream"]).To(HaveLen(1))
				Expect(matrix.m["upstream"][Vertex{}]).To(HaveLen(2))
			})
		})

		When("child vertices is empty", func() {
			It("adds an edge", func() {
				matrix.AddEdges("upstream", Vertex{"a", "b", "c", "d"}, []Vertex{})

				Expect(matrix.m).To(HaveLen(1))
				Expect(matrix.m["upstream"]).To(HaveLen(1))
				Expect(matrix.m["upstream"][Vertex{"a", "b", "c", "d"}]).To(HaveLen(0))
			})
		})
	})

	Context("get", func() {

		When("matrix is empty", func() {
			When("named edge is empty", func() {
				It("returns nil", func() {
					v := matrix.Get("", Vertex{"a", "b", "c", "d"})

					Expect(v).To(BeNil())
				})
			})

			When("parent is empty", func() {
				It("returns nil", func() {
					v := matrix.Get("upstream", Vertex{})

					Expect(v).To(BeNil())
				})
			})

			When("parent is not empty", func() {
				It("returns nil", func() {
					v := matrix.Get("upstream", Vertex{"a", "b", "c", "d"})

					Expect(v).To(BeNil())
				})
			})
		})

		When("matrix is not empty", func() {
			BeforeEach(func() {
				matrix.AddEdges("upstream", Vertex{"a", "b", "c", "d"}, []Vertex{{"e", "f", "g", "h"}, {"i", "j", "k", "l"}})
			})
			When("named edge is empty", func() {
				It("returns nil", func() {
					v := matrix.Get("", Vertex{"a", "b", "c", "d"})

					Expect(v).To(BeNil())
				})
			})

			When("parent is empty", func() {
				It("returns nil", func() {
					v := matrix.Get("upstream", Vertex{})

					Expect(v).To(BeNil())
				})
			})

			When("parent is not empty", func() {
				It("returns nil", func() {
					v := matrix.Get("upstream", Vertex{"a", "b", "c", "d"})

					Expect(v).To(Equal([]Vertex{{"e", "f", "g", "h"}, {"i", "j", "k", "l"}}))
				})
			})
		})
	})

	Context("number of edges", func() {

		When("matrix is empty", func() {
			When("edge name is empty", func() {
				It("returns zero", func() {
					n := matrix.NumberOfEdges("")

					Expect(n).To(Equal(0))
				})
			})

			When("edge name is empty", func() {
				It("returns zero", func() {
					n := matrix.NumberOfEdges("upstream")

					Expect(n).To(Equal(0))
				})
			})
		})

		When("matrix is not empty", func() {
			BeforeEach(func() {
				matrix.AddEdges("upstream", Vertex{"a", "b", "c", "d"}, []Vertex{{"e", "f", "g", "h"}, {"i", "j", "k", "l"}})
			})

			When("edge name is empty", func() {
				It("returns zero", func() {
					n := matrix.NumberOfEdges("")

					Expect(n).To(Equal(0))
				})
			})

			When("edge name is not empty", func() {
				It("returns nil", func() {
					n := matrix.NumberOfEdges("upstream")

					Expect(n).To(Equal(1))
				})
			})
		})
	})
})
