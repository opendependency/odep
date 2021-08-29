package v1

import (
	"errors"
	"fmt"
	"regexp"
)

var isLowercaseAlphanumericDashDot = regexp.MustCompile(`^[a-z0-9-.]+$`).MatchString

// Validate checks if the specification constraints are fulfilled.
func (x *Module) Validate() error {
	if err := validateModuleNamespace(x.Namespace); err != nil {
		return fmt.Errorf("namespace: %w", err)
	}
	if err := validateModuleName(x.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if err := validateModuleType(x.Type); err != nil {
		return fmt.Errorf("type: %w", err)
	}

	if err := validateModuleVersion(x.Version); err != nil {
		return fmt.Errorf("version: %w", err)
	}

	if err := validateModuleAnnotations(x.Annotations); err != nil {
		return fmt.Errorf("annotations: %w", err)
	}

	if err := validateModuleDependencies(x.Dependencies); err != nil {
		return fmt.Errorf("dependencies: %w", err)
	}

	return nil
}

func validateModuleNamespace(namespace string) error {
	return mustFulfilConstraints(
		func() error {
			return mustHaveMinMaxLength(namespace, 1, 63)
		},
		func() error {
			return mustBeLowercaseAlphanumericDashDot(namespace)
		},
		func() error {
			return mustStartWithLowercaseAlphabeticCharacter(namespace)
		},
		func() error {
			return mustEndWithLowercaseAlphanumericCharacter(namespace)
		},
	)
}

func validateModuleName(name string) error {
	return mustFulfilConstraints(
		func() error {
			return mustHaveMinMaxLength(name, 1, 63)
		},
		func() error {
			return mustBeLowercaseAlphanumericDashDot(name)
		},
		func() error {
			return mustStartWithLowercaseAlphabeticCharacter(name)
		},
		func() error {
			return mustEndWithLowercaseAlphanumericCharacter(name)
		},
	)
}

func validateModuleType(type_ string) error {
	return mustFulfilConstraints(
		func() error {
			return mustHaveMinMaxLength(type_, 1, 63)
		},
		func() error {
			return mustBeLowercaseAlphanumericDashDot(type_)
		},
		func() error {
			return mustStartWithLowercaseAlphabeticCharacter(type_)
		},
		func() error {
			return mustEndWithLowercaseAlphanumericCharacter(type_)
		},
	)
}

func validateModuleVersion(moduleVersion *ModuleVersion) error {
	if moduleVersion == nil {
		return errors.New("must be set")
	}

	return moduleVersion.Validate()
}

// Validate checks if the specification constraints are fulfilled.
func (x *ModuleVersion) Validate() error {
	if err := validateModuleVersionName(x.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if x.Schema != nil {
		if err := validateModuleVersionSchema(*x.Schema); err != nil {
			return fmt.Errorf("schema: %w", err)
		}
	}

	for i, v := range x.Replaces {
		if err := validateModuleVersionName(v); err != nil {
			return fmt.Errorf("replaces: index %d: %w", i, err)
		}
	}

	return nil
}

func validateModuleVersionName(name string) error {
	return mustFulfilConstraints(
		func() error {
			return mustHaveMinMaxLength(name, 1, 63)
		},
		func() error {
			return mustBeLowercaseAlphanumericDashDot(name)
		},
		func() error {
			return mustStartWithLowercaseAlphanumericCharacter(name)
		},
		func() error {
			return mustEndWithLowercaseAlphanumericCharacter(name)
		},
	)
}

func validateModuleVersionSchema(schema string) error {
	return mustFulfilConstraints(
		func() error {
			return mustHaveMinMaxLength(schema, 1, 63)
		},
		func() error {
			return mustBeLowercaseAlphanumericDashDot(schema)
		},
		func() error {
			return mustStartWithLowercaseAlphabeticCharacter(schema)
		},
		func() error {
			return mustEndWithLowercaseAlphanumericCharacter(schema)
		},
	)
}

func validateModuleAnnotations(annotations map[string]string) error {
	if annotations == nil || len(annotations) == 0 {
		return nil
	}

	for k, v := range annotations {
		if err := validateModuleAnnotationKey(k); err != nil {
			return fmt.Errorf("key %q: %w", k, err)
		}
		if err := validateModuleAnnotationValue(v); err != nil {
			return fmt.Errorf("value of key %q: %w", k, err)
		}
	}

	return nil
}

func validateModuleAnnotationKey(key string) error {
	return mustFulfilConstraints(
		func() error {
			return mustHaveMinMaxLength(key, 1, 63)
		},
		func() error {
			return mustBeLowercaseAlphanumericDashDot(key)
		},
		func() error {
			return mustStartWithLowercaseAlphabeticCharacter(key)
		},
		func() error {
			return mustEndWithLowercaseAlphanumericCharacter(key)
		},
	)
}

func validateModuleAnnotationValue(value string) error {
	return mustFulfilConstraints(
		func() error {
			return mustHaveMinMaxLength(value, 0, 253)
		},
	)
}

func validateModuleDependencies(moduleDependencies []*ModuleDependency) error {
	if moduleDependencies == nil || len(moduleDependencies) == 0 {
		return nil
	}

	for i := 0; i < len(moduleDependencies); i++ {
		moduleDependency := moduleDependencies[i]
		if err := moduleDependency.Validate(); err != nil {
			return fmt.Errorf("index %d: %w", i, err)
		}
	}

	return nil
}

// Validate checks if the specification constraints are fulfilled.
func (x *ModuleDependency) Validate() error {
	if err := validateModuleNamespace(x.Namespace); err != nil {
		return fmt.Errorf("namespace: %w", err)
	}
	if err := validateModuleName(x.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if err := validateModuleType(x.Type); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	if err := validateModuleVersionName(x.Version); err != nil {
		return fmt.Errorf("version: %w", err)
	}

	return nil
}

func mustFulfilConstraints(constraints ...func() error) error {
	for _, constraint := range constraints {
		if err := constraint(); err != nil {
			return err
		}
	}
	return nil
}

func mustHaveMinMaxLength(value string, minLen int, maxLen int) error {
	if minLen < 0 {
		return fmt.Errorf("min length must be greater or equal 0")
	}
	if maxLen < 0 {
		return fmt.Errorf("max length must be greater or equal 0")
	}
	if minLen > maxLen {
		return fmt.Errorf("min length must be less or equal max length")
	}

	l := len(value)

	if l < minLen {
		return fmt.Errorf("must have at least %d characters", minLen)
	}
	if l > maxLen {
		return fmt.Errorf("must have at most %d characters", maxLen)
	}

	return nil
}

func mustBeLowercaseAlphanumericDashDot(value string) error {
	if len(value) == 0 {
		return nil
	}

	if !isLowercaseAlphanumericDashDot(value) {
		return fmt.Errorf("must contain only lowercase alphanumeric characters, '-' or '.'")
	}

	return nil
}

func mustStartWithLowercaseAlphabeticCharacter(value string) error {
	if len(value) < 1 {
		return nil
	}

	firstCharacter := rune(value[0])

	if firstCharacter >= 'a' && firstCharacter <= 'z' {
		return nil
	}

	return fmt.Errorf("must start with lowercase alphabetic character")
}

func mustStartWithLowercaseAlphanumericCharacter(value string) error {
	if len(value) == 0 {
		return nil
	}

	firstCharacter := rune(value[0])

	if (firstCharacter >= 'a' && firstCharacter <= 'z') || (firstCharacter >= '0' && firstCharacter <= '9') {
		return nil
	}

	return fmt.Errorf("must start with lowercase alphanumeric character")
}

func mustEndWithLowercaseAlphanumericCharacter(value string) error {
	if len(value) < 1 {
		return nil
	}

	lastCharacter := rune(value[len(value)-1])

	if (lastCharacter >= 'a' && lastCharacter <= 'z') || (lastCharacter >= '0' && lastCharacter <= '9') {
		return nil
	}

	return fmt.Errorf("must end with lowercase alphanumeric character")
}
