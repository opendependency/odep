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
	"errors"
	"fmt"
	"sync"

	spec "github.com/opendependency/go-spec/pkg/spec/v1"
	"google.golang.org/protobuf/proto"
)

// NewInMemoryRepository creates a new in-memory repository.
func NewInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{
		data: map[string]map[string]map[string]map[string]*spec.Module{},
	}
}

var _ Repository = (*inMemoryRepository)(nil)

type inMemoryRepository struct {
	mux  sync.RWMutex
	data map[string]map[string]map[string]map[string]*spec.Module
}

func (r *inMemoryRepository) AddModule(module *spec.Module) error {
	if module == nil {
		return errors.New("module must not be nil")
	}

	if err := module.Validate(); err != nil {
		return fmt.Errorf("module validation failed: %w", err)
	}

	clone := proto.Clone(module).(*spec.Module)

	r.mux.Lock()

	moduleNames := r.data[clone.Namespace]
	if moduleNames == nil {
		moduleNames = map[string]map[string]map[string]*spec.Module{}
		r.data[clone.Namespace] = moduleNames
	}

	moduleTypes := moduleNames[clone.Name]
	if moduleTypes == nil {
		moduleTypes = map[string]map[string]*spec.Module{}
		moduleNames[clone.Name] = moduleTypes
	}

	moduleVersions := moduleTypes[clone.Type]
	if moduleVersions == nil {
		moduleVersions = map[string]*spec.Module{}
		moduleTypes[clone.Type] = moduleVersions
	}

	moduleVersions[clone.Version.Name] = clone

	r.mux.Unlock()

	return nil
}

func (r *inMemoryRepository) DeleteNamespace(namespace string) error {
	r.mux.Lock()
	delete(r.data, namespace)
	r.mux.Unlock()

	return nil
}

func (r *inMemoryRepository) DeleteModule(namespace string, name string) error {
	r.mux.Lock()
	moduleNames := r.data[namespace]
	if moduleNames != nil {
		delete(moduleNames, name)
	}
	r.mux.Unlock()

	return nil
}

func (r *inMemoryRepository) DeleteModuleType(namespace string, name string, type_ string) error {
	r.mux.Lock()
	if moduleNames := r.data[namespace]; moduleNames != nil {
		if moduleTypes := moduleNames[name]; moduleTypes != nil {
			delete(moduleTypes, type_)
		}
	}
	r.mux.Unlock()

	return nil
}

func (r *inMemoryRepository) DeleteModuleVersion(namespace string, name string, type_ string, version string) error {
	r.mux.Lock()
	if moduleNames := r.data[namespace]; moduleNames != nil {
		if moduleTypes := moduleNames[name]; moduleTypes != nil {
			if moduleVersions := moduleTypes[type_]; moduleVersions != nil {
				delete(moduleVersions, version)
			}
		}
	}
	r.mux.Unlock()

	return nil
}

func (r *inMemoryRepository) GetModule(namespace string, name string, type_ string, version string) (*spec.Module, error) {
	var module *spec.Module

	r.mux.RLock()
	if moduleNames := r.data[namespace]; moduleNames != nil {
		if moduleTypes := moduleNames[name]; moduleTypes != nil {
			if moduleVersions := moduleTypes[type_]; moduleVersions != nil {
				if m, ok := moduleVersions[version]; ok {
					module = proto.Clone(m).(*spec.Module)
				}
			}
		}
	}
	r.mux.RUnlock()

	if module != nil {
		return module, nil
	}

	return nil, fmt.Errorf("not found")
}

func (r *inMemoryRepository) ListModuleNamespaces() ([]string, error) {
	var namespaces []string

	r.mux.RLock()
	for k := range r.data {
		namespaces = append(namespaces, k)
	}
	r.mux.RUnlock()

	return namespaces, nil
}

func (r *inMemoryRepository) ListModuleNames(namespace string) ([]string, error) {
	var names []string

	r.mux.RLock()
	for k := range r.data[namespace] {
		names = append(names, k)
	}
	r.mux.RUnlock()

	return names, nil
}

func (r *inMemoryRepository) ListModuleTypes(namespace string, name string) ([]string, error) {
	var types []string

	r.mux.RLock()
	if moduleNames := r.data[namespace]; moduleNames != nil {
		for k := range moduleNames[name] {
			types = append(types, k)
		}
	}
	r.mux.RUnlock()

	return types, nil
}

func (r *inMemoryRepository) ListModuleVersions(namespace string, name string, type_ string) ([]string, error) {
	var versions []string

	r.mux.RLock()
	if moduleNames := r.data[namespace]; moduleNames != nil {
		if moduleTypes := moduleNames[name]; moduleTypes != nil {
			for k := range moduleTypes[type_] {
				versions = append(versions, k)
			}
		}
	}
	r.mux.RUnlock()

	return versions, nil
}
