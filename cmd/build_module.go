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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	specv1 "github.com/opendependency/go-spec/pkg/spec/v1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"sigs.k8s.io/yaml"
)

// NewBuildModuleCommand creates a new build module command.
func NewBuildModuleCommand() *cobra.Command {
	var (
		module = &specv1.Module{}

		moduleFile string

		moduleNamespace string
		moduleName      string
		moduleType      string

		moduleVersionName     string
		moduleVersionSchema   string
		moduleVersionReplaces []string

		moduleAnnotations            map[string]string
		moduleUpstreamDependencies   []string
		moduleDownstreamDependencies []string

		moduleOutput       string
		moduleOutputPretty bool
	)

	buildModuleCmd := &cobra.Command{
		Use:          "module",
		Short:        "Builds a module.",
		Long:         `Builds a module.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if moduleFile != "" {
				if _, err := os.Stat(moduleFile); os.IsNotExist(err) {
					return fmt.Errorf("module file does not exist")
				}
				fileContent, err := ioutil.ReadFile(moduleFile)
				if err != nil {
					return fmt.Errorf("could not read module file: %w", err)
				}

				switch filepath.Ext(moduleFile) {
				case ".json":
					if err := protojson.Unmarshal(fileContent, module); err != nil {
						return fmt.Errorf("could not unmarshal json: %w", err)
					}
				case ".yaml":
					json, err := yaml.YAMLToJSON(fileContent)
					if err != nil {
						return fmt.Errorf("could not convert yaml to json: %w", err)
					}
					if err := protojson.Unmarshal(json, module); err != nil {
						return fmt.Errorf("could not unmarshal json: %w", err)
					}
				}
			}

			if moduleNamespace != "" {
				module.Namespace = moduleNamespace
			}
			if moduleName != "" {
				module.Name = moduleName
			}
			if moduleType != "" {
				module.Type = moduleType
			}

			if module.Version == nil {
				module.Version = &specv1.ModuleVersion{}
			}
			if moduleVersionName != "" {
				module.Version.Name = moduleVersionName
			}
			if moduleVersionSchema != "" {
				module.Version.Schema = &moduleVersionSchema
			}
			if len(moduleVersionReplaces) > 0 {
				module.Version.Replaces = moduleVersionReplaces
			}

			if moduleAnnotations != nil && len(moduleAnnotations) > 0 {
				module.Annotations = moduleAnnotations
			}

			upstream := specv1.DependencyDirection_UPSTREAM
			downstream := specv1.DependencyDirection_DOWNSTREAM

			for _, dependency := range moduleUpstreamDependencies {
				moduleDependency, err := parseModuleDependency(dependency)
				if err != nil {
					return fmt.Errorf("could not parse upstream module dependency %q: %w", dependency, err)
				}
				moduleDependency.Direction = &upstream
				module.Dependencies = append(module.Dependencies, moduleDependency)
			}

			for _, dependency := range moduleDownstreamDependencies {
				moduleDependency, err := parseModuleDependency(dependency)
				if err != nil {
					return fmt.Errorf("could not parse downstream module dependency %q: %w", dependency, err)
				}
				moduleDependency.Direction = &downstream
				module.Dependencies = append(module.Dependencies, moduleDependency)
			}

			if err := module.Validate(); err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			switch moduleOutput {
			case "json":
				marshalJson, err := protojson.Marshal(module)
				if err != nil {
					return fmt.Errorf("could not marshal to json: %w", err)
				}

				buf := &bytes.Buffer{}
				if moduleOutputPretty {
					if err := json.Indent(buf, marshalJson, "", "  "); err != nil {
						return fmt.Errorf("could not indent json: %w", err)
					}
				} else {
					if err := json.Compact(buf, marshalJson); err != nil {
						return fmt.Errorf("could not compact json: %w", err)
					}
				}

				cmd.Print(buf.String())
			case "yaml":
				marshalJson, err := protojson.Marshal(module)
				if err != nil {
					return fmt.Errorf("could not marshal to json: %w", err)
				}
				marshalledYaml, err := yaml.JSONToYAML(marshalJson)
				if err != nil {
					return fmt.Errorf("could not convert to yaml: %w", err)
				}
				cmd.Print(string(marshalledYaml))
			default:
				cmd.Printf("Module %s %s %s %s built.\n", module.Namespace, module.Name, module.Type, module.Version.Name)
			}

			return nil
		},
	}

	buildModuleCmd.Flags().StringVarP(&moduleFile, "from-file", "f", "", "From-file specifies a module file. Supported formats: json, yaml")

	buildModuleCmd.Flags().StringVar(&moduleNamespace, "namespace", "", "Namespace defines the module namespace.")
	buildModuleCmd.Flags().StringVar(&moduleName, "name", "", "Name defines the module name.")
	buildModuleCmd.Flags().StringVar(&moduleType, "type", "", "Type defines the module type.")

	buildModuleCmd.Flags().StringVar(&moduleVersionName, "version-name", "", "Version name defines the module version name.")
	buildModuleCmd.Flags().StringVar(&moduleVersionSchema, "version-schema", "", "Version schema defines the module version schema.")
	buildModuleCmd.Flags().StringSliceVar(&moduleVersionReplaces, "version-replaces", nil, "Version replaces defines a list of previous module version replaced by this version.")

	buildModuleCmd.Flags().StringToStringVar(&moduleAnnotations, "annotations", nil, "Annotations defines arbitrary module metadata.")
	buildModuleCmd.Flags().StringSliceVar(&moduleUpstreamDependencies, "upstream-dependencies", nil, "Upstream dependencies specifies all upstream dependencies. Notation: <NAMESPACE>:<NAME>:<TYPE>:<VERSION>")
	buildModuleCmd.Flags().StringSliceVar(&moduleDownstreamDependencies, "downstream-dependencies", nil, "Downstream dependencies specifies all downstream dependencies. Notation: <NAMESPACE>:<NAME>:<TYPE>:<VERSION>")

	buildModuleCmd.Flags().StringVarP(&moduleOutput, "output", "o", "", "Output defines the output format. Supported formats: json, yaml")
	buildModuleCmd.Flags().BoolVarP(&moduleOutputPretty, "pretty", "", false, "Pretty prints the output in multiple lines with indents.")

	return buildModuleCmd
}

func parseModuleDependency(dependency string) (*specv1.ModuleDependency, error) {
	dependencyParts := strings.Split(dependency, ":")

	if len(dependencyParts) != 4 {
		return nil, fmt.Errorf("must be of notation '<NAMESPACE>:<NAME>:<TYPE>:<VERSION>'")
	}

	return &specv1.ModuleDependency{
		Namespace: dependencyParts[0],
		Name:      dependencyParts[1],
		Type:      dependencyParts[2],
		Version:   dependencyParts[3],
	}, nil
}
