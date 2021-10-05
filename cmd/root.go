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
	"github.com/spf13/cobra"
)

// NewRootCommand creates a new root command.
func NewRootCommand(ctx Context) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "odep",
		Short: "odep manages OpenDependency modules.",
		Long:  `odep manages OpenDependency modules.`,
		// see https://github.com/spf13/cobra/issues/706#issuecomment-488340260
		Args: cobra.NoArgs,
	}

	rootCmd.AddCommand(NewBuildCommand(ctx))

	return rootCmd
}
