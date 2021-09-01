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
	"container/list"
	"errors"
	"fmt"

	spec "github.com/opendependency/go-spec/pkg/spec/v1"
)

// Vertex represents a module within a graph.
type Vertex struct {
	Namespace string
	Name      string
	Type      string
	Version   string
}

func (v *Vertex) String() string {
	return fmt.Sprintf("%s:%s:%s:%s", v.Namespace, v.Name, v.Type, v.Version)
}

// Graph represents a module graph containing all edges to other modules.
type Graph interface {
	// AddModule adds the given module.
	AddModule(module *spec.Module) error
	// TraverseDependOnEdgesBFS begins at vertex s and traverse over all depend-on edges
	// using breadth-first search.
	// The given function fn is called for each vertex and its direct depend-on edge vertices.
	// The function fn returning true continues the traversal while returning false stops the traversal.
	// The first function fn call has vertex s as parent p.
	TraverseDependOnEdgesBFS(s Vertex, fn func(p Vertex, v []Vertex) bool)
	// TraverseDependOnEdgesDFS begins at Vertex s and traverse over all depend-on edges
	// using depth-first search.
	// The given function fn is called for each vertex and its depend-on edge vertices.
	// The function fn returning true continues the traversal while returning false stops the traversal.
	// The first function fn call has an empty vertex as parent p.
	TraverseDependOnEdgesDFS(s Vertex, fn func(p Vertex, v Vertex) bool)
	// TraverseUsedByEdgesBFS begins at vertex s and traverse over all used-by edges
	// using breadth-first search.
	// The given function fn is called for each vertex and its direct used-by edge vertices.
	// The function fn returning true continues the traversal while returning false stops the traversal.
	// The first function fn call has vertex s as parent p.
	TraverseUsedByEdgesBFS(s Vertex, fn func(p Vertex, v []Vertex) bool)
	// TraverseUsedByEdgesDFS begins at Vertex s and traverse over all used-by edges
	// using depth-first search.
	// The given function fn is called for each vertex and its used-by edge vertices.
	// The function fn returning true continues the traversal while returning false stops the traversal.
	// The first function fn call has an empty vertex as parent p.
	TraverseUsedByEdgesDFS(s Vertex, fn func(p Vertex, v Vertex) bool)
	// TraverseRequiredForEdgesBFS begins at vertex s and traverse over all required-for edges
	// using breadth-first search.
	// The given function fn is called for each vertex and its direct required-for edge vertices.
	// The function fn returning true continues the traversal while returning false stops the traversal.
	// The first function fn call has vertex s as parent p.
	TraverseRequiredForEdgesBFS(s Vertex, fn func(p Vertex, v []Vertex) bool)
	// TraverseRequiredForEdgesDFS begins at Vertex s and traverse over all required-for edges
	// using depth-first search.
	// The given function fn is called for each vertex and its required-for edge vertices.
	// The function fn returning true continues the traversal while returning false stops the traversal.
	// The first function fn call has an empty vertex as parent p.
	TraverseRequiredForEdgesDFS(s Vertex, fn func(p Vertex, v Vertex) bool)
	// TraverseRequireEdgesBFS begins at vertex s and traverse over all require edges
	// using breadth-first search.
	// The given function fn is called for each vertex and its direct require edge vertices.
	// The function fn returning true continues the traversal while returning false stops the traversal.
	// The first function fn call has vertex s as parent p.
	TraverseRequireEdgesBFS(s Vertex, fn func(p Vertex, v []Vertex) bool)
	// TraverseRequireEdgesDFS begins at Vertex s and traverse over all require edges
	// using depth-first search.
	// The given function fn is called for each vertex and its require edge vertices.
	// The function fn returning true continues the traversal while returning false stops the traversal.
	// The first function fn call has an empty vertex as parent p.
	TraverseRequireEdgesDFS(s Vertex, fn func(p Vertex, v Vertex) bool)
}

const (
	// dependsOnEdge represents edges where vertex A depend on vertex B.
	// Opposite: vertex B is used by vertex A.
	dependsOnEdge = "depends-on"
	// usedByEdge represents edges where vertex A is used by vertex B.
	// Opposite: vertex B depends on vertex A.
	usedByEdge = "used-by"
	// requiredForEdge represents edges where vertex A is required for vertex B.
	// Opposite: vertex B requires vertex A.
	requiredForEdge = "required-for"
	// requireEdge represents edges where vertex A requires vertex B.
	// Opposite: vertex B is required for vertex A.
	requireEdge = "require"
)

// NewGraph creates a new graph with the given AdjacentMatrix as underlying matrix.
func NewGraph(m AdjacentMatrix) *graph {
	return &graph{
		m: m,
	}
}

var _ Graph = (*graph)(nil)

type graph struct {
	m AdjacentMatrix
}

func (g *graph) AddModule(module *spec.Module) error {
	if module == nil {
		return errors.New("module must not be nil")
	}

	if err := module.Validate(); err != nil {
		return fmt.Errorf("module validation failed: %w", err)
	}

	p := Vertex{
		Namespace: module.Namespace,
		Name:      module.Name,
		Type:      module.Type,
		Version:   module.Version.Name,
	}

	for _, dependency := range module.Dependencies {
		v := Vertex{
			Namespace: dependency.Namespace,
			Name:      dependency.Name,
			Type:      dependency.Type,
			Version:   dependency.Version,
		}

		if dependency.Direction == nil || *dependency.Direction == spec.DependencyDirection_UPSTREAM {
			g.m.AddEdge(dependsOnEdge, p, v)
			g.m.AddEdge(usedByEdge, v, p)
		} else {
			g.m.AddEdge(requiredForEdge, p, v)
			g.m.AddEdge(requireEdge, v, p)
		}
	}

	return nil
}

func (g *graph) TraverseDependOnEdgesBFS(s Vertex, fn func(p Vertex, v []Vertex) bool) {
	g.traverseBFS(dependsOnEdge, s, fn)
}

func (g *graph) TraverseDependOnEdgesDFS(s Vertex, fn func(p Vertex, v Vertex) bool) {
	g.traverseDFS(dependsOnEdge, s, fn)
}

func (g *graph) TraverseUsedByEdgesBFS(s Vertex, fn func(p Vertex, v []Vertex) bool) {
	g.traverseBFS(usedByEdge, s, fn)
}

func (g *graph) TraverseUsedByEdgesDFS(s Vertex, fn func(p Vertex, v Vertex) bool) {
	g.traverseDFS(usedByEdge, s, fn)
}

func (g *graph) TraverseRequiredForEdgesBFS(s Vertex, fn func(p Vertex, v []Vertex) bool) {
	g.traverseBFS(requiredForEdge, s, fn)
}

func (g *graph) TraverseRequiredForEdgesDFS(s Vertex, fn func(p Vertex, v Vertex) bool) {
	g.traverseDFS(requiredForEdge, s, fn)
}

func (g *graph) TraverseRequireEdgesBFS(s Vertex, fn func(p Vertex, v []Vertex) bool) {
	g.traverseBFS(requireEdge, s, fn)
}

func (g *graph) TraverseRequireEdgesDFS(s Vertex, fn func(p Vertex, v Vertex) bool) {
	g.traverseDFS(requireEdge, s, fn)
}

func (g *graph) traverseBFS(edgeName string, s Vertex, fn func(p Vertex, v []Vertex) bool) {
	// track visited vertices
	visited := map[Vertex]bool{}
	// track vertices to visit
	queue := list.New()
	queue.PushBack(s)
	// mark start vertex as visited
	visited[s] = true

	for queue.Len() > 0 {
		qv := queue.Front()

		// iterate through all children
		children := g.m.Get(edgeName, qv.Value.(Vertex))

		if ok := fn(qv.Value.(Vertex), children); !ok {
			return
		}

		for _, child := range children {
			if ok := visited[child]; !ok {
				visited[child] = true
				queue.PushBack(child)
			}
		}

		queue.Remove(qv)
	}
}

func (g *graph) traverseDFS(edgeName string, s Vertex, fn func(p Vertex, v Vertex) bool) {
	var emptyVertex Vertex

	// track visited vertices
	visited := map[Vertex]bool{}

	stack := &vertexPairStack{}
	stack.Push(emptyVertex, s)

	for {
		p, v, err := stack.Pop()
		if err == emptyStackErr {
			return
		}

		// mark as visited
		visited[v] = true

		if ok := fn(p, v); !ok {
			return
		}

		// add all children
		children := g.m.Get(edgeName, v)
		for _, child := range children {
			if ok := visited[child]; !ok {
				stack.Push(v, child)
			}
		}
	}
}

var emptyStackErr = errors.New("empty stack")

type vertexPair struct {
	k Vertex
	v Vertex
}

type vertexPairStack struct {
	s []vertexPair
}

func (s *vertexPairStack) Push(k Vertex, v Vertex) {
	s.s = append(s.s, vertexPair{k, v})
}

func (s *vertexPairStack) Pop() (Vertex, Vertex, error) {
	l := len(s.s)
	if l == 0 {
		return Vertex{}, Vertex{}, emptyStackErr
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res.k, res.v, nil
}
