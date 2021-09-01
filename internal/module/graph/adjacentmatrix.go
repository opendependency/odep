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

// AdjacentMatrix represents a directed graph through an adjacent matrix.
type AdjacentMatrix interface {
	// AddEdge adds a named edge between vertex p and vertex c.
	AddEdge(name string, p Vertex, c Vertex)
	// AddEdges adds a named edge between vertex p and vertices c.
	AddEdges(name string, p Vertex, c []Vertex)
	// Get gets all vertices of a named edge on vertex v.
	Get(name string, v Vertex) []Vertex
	// NumberOfEdges gets the number of named edges.
	NumberOfEdges(name string) int
}
