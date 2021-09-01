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
	"sync"
)

// NewInMemoryAdjacentMatrix creates a new in-memory adjacent matrix.
func NewInMemoryAdjacentMatrix() *inMemoryAdjacentMatrix {
	return &inMemoryAdjacentMatrix{
		m: map[string]map[Vertex][]Vertex{},
	}
}

var _ AdjacentMatrix = (*inMemoryAdjacentMatrix)(nil)

type inMemoryAdjacentMatrix struct {
	mux sync.RWMutex
	m   map[string]map[Vertex][]Vertex
}

func (a *inMemoryAdjacentMatrix) AddEdge(name string, p Vertex, c Vertex) {
	a.mux.Lock()
	matrix, ok := a.m[name]
	if !ok {
		matrix = map[Vertex][]Vertex{}
		a.m[name] = matrix
	}
	matrix[p] = append(matrix[p], c)
	a.mux.Unlock()
}

func (a *inMemoryAdjacentMatrix) AddEdges(name string, p Vertex, c []Vertex) {
	a.mux.Lock()
	matrix, ok := a.m[name]
	if !ok {
		matrix = map[Vertex][]Vertex{}
		a.m[name] = matrix
	}
	matrix[p] = append(matrix[p], c...)
	a.mux.Unlock()
}

func (a *inMemoryAdjacentMatrix) Get(name string, v Vertex) []Vertex {
	a.mux.RLock()
	defer a.mux.RUnlock()
	matrix, ok := a.m[name]
	if !ok {
		return nil
	}
	return matrix[v]
}

func (a *inMemoryAdjacentMatrix) NumberOfEdges(name string) int {
	return len(a.m[name])
}
