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
	"github.com/opendependency/odep/internal/module/repository"
)

// Context represents the command context.
type Context interface {
	// ModuleRepository provides the module repository to interact with modules.
	ModuleRepository() repository.Repository
}

// NewContext creates a new command context.
func NewContext(moduleRepository repository.Repository) *context {
	return &context{
		moduleRepository: moduleRepository,
	}
}

var _ Context = &context{}

type context struct {
	moduleRepository repository.Repository
}

func (c *context) ModuleRepository() repository.Repository {
	return c.moduleRepository
}
