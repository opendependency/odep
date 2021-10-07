package cmd

import (
	"github.com/opendependency/odep/internal/module/repository"
	"github.com/spf13/cobra"
)

// Providers provide various resources.
type Providers struct {
	// ModuleRepository provides a module repository.
	ModuleRepository func(*cobra.Command) ModuleRepositoryProvider
}

// NewDefaultProviders creates the default providers.
func NewDefaultProviders() *Providers {
	return &Providers{
		ModuleRepository: NewDefaultModuleRepositoryProvider(),
	}
}

// ModuleRepositoryProvider provides a module repository.
type ModuleRepositoryProvider func() repository.Repository

func NewDefaultModuleRepositoryProvider() func(*cobra.Command) ModuleRepositoryProvider {
	return func(targetCmd *cobra.Command) ModuleRepositoryProvider {
		var(

		)
		targetCmd.Flags().StringVar(&moduleType, "type", "", "Type defines the module type.")
	}
}
