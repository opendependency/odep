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
	spec "github.com/opendependency/go-spec/pkg/spec/v1"
)

// Repository provides access to modules stored in a backend.
type Repository interface {
	// AddModule adds the given module.
	AddModule(module *spec.Module) error
	// DeleteNamespace deletes a whole module namespace with all modules.
	DeleteNamespace(namespace string) error
	// DeleteModule deletes a specific module.
	DeleteModule(namespace string, name string) error
	// DeleteModuleType deletes a specific module type.
	DeleteModuleType(namespace string, name string, type_ string) error
	// DeleteModuleVersion deletes a specific module version.
	DeleteModuleVersion(namespace string, name string, type_ string, version string) error
	// GetModule gets a specific module.
	GetModule(namespace string, name string, type_ string, version string) (*spec.Module, error)
	// ListModuleNamespaces list all module namespaces.
	ListModuleNamespaces() ([]string, error)
	// ListModuleNames list all module names within a namespace.
	ListModuleNames(namespace string) ([]string, error)
	// ListModuleTypes list all module types of a module.
	ListModuleTypes(namespace string, name string) ([]string, error)
	// ListModuleVersions list all module versions of a module.
	ListModuleVersions(namespace string, name string, type_ string) ([]string, error)
}
