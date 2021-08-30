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
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/flock"
	spec "github.com/opendependency/go-spec/pkg/spec/v1"
	"google.golang.org/protobuf/proto"
)

const (
	modulesDirectory    = "modules"
	moduleFileExtension = "module.bin"
)

// NewFileRepository creates a new file repository under the given path.
func NewFileRepository(path string) (*fileRepository, error) {
	absDir, err := filepath.Abs(filepath.Join(path, modulesDirectory))
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path: %w", err)
	}

	if err := os.MkdirAll(absDir, os.ModePerm); err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("could not create directory: %w", err)
	}

	return &fileRepository{
		path: absDir,
	}, nil
}

var _ Repository = (*fileRepository)(nil)

type fileRepository struct {
	path string
}

func (r *fileRepository) AddModule(module *spec.Module) (rerr error) {
	if module == nil {
		return errors.New("module must not be nil")
	}

	if err := module.Validate(); err != nil {
		return fmt.Errorf("module validation failed: %w", err)
	}

	serializedModule, err := proto.Marshal(module)
	if err != nil {
		return fmt.Errorf("could not marhsal proto: %w", err)
	}

	if err := os.MkdirAll(r.getAbsoluteModuleTypeDirectoryPath(module.Namespace, module.Name, module.Type), os.ModePerm); err != nil && !os.IsExist(err) {
		return fmt.Errorf("could not create directory: %w", err)
	}

	targetAbsModuleFilePath := r.getAbsoluteModuleFilePath(module.Namespace, module.Name, module.Type, module.Version.Name)

	l := r.newFileLock(targetAbsModuleFilePath)
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	locked, err := l.TryLockContext(lockCtx, 500*time.Millisecond)
	if !locked || err != nil {
		return fmt.Errorf("could not lock: %s", l.Path())
	}

	defer func() {
		if err := l.Unlock(); err != nil {
			if rerr != nil {
				rerr = fmt.Errorf("%s ; could not unlock: %w", rerr.Error(), err)
			}
			rerr = fmt.Errorf("could not unlock: %w", err)
		}
	}()

	if err := ioutil.WriteFile(targetAbsModuleFilePath, serializedModule, os.ModePerm); err != nil {
		return fmt.Errorf("could not write module file: %w", err)
	}

	return nil
}

func (r *fileRepository) newFileLock(absFilePath string) *flock.Flock {
	return flock.New(absFilePath + ".lock")
}

func (r *fileRepository) getAbsoluteModuleNamespaceDirectoryPath(namespace string) string {
	return path.Join(r.path, namespace)
}

func (r *fileRepository) getAbsoluteModuleNameDirectoryPath(namespace string, name string) string {
	return path.Join(r.path, namespace, name)
}

func (r *fileRepository) getAbsoluteModuleTypeDirectoryPath(namespace string, name string, type_ string) string {
	return path.Join(r.path, namespace, name, type_)
}

func (r *fileRepository) getAbsoluteModuleFilePath(namespace string, name string, type_ string, version string) string {
	return path.Join(r.path, namespace, name, type_, fmt.Sprintf("%s.%s", version, moduleFileExtension))
}

func (r *fileRepository) DeleteNamespace(namespace string) error {
	if err := os.RemoveAll(r.getAbsoluteModuleNamespaceDirectoryPath(namespace)); err != nil {
		return err
	}
	return nil
}

func (r *fileRepository) DeleteModule(namespace string, name string) error {
	if err := os.RemoveAll(r.getAbsoluteModuleNameDirectoryPath(namespace, name)); err != nil {
		return err
	}
	return r.cleanup(r.getAbsoluteModuleNamespaceDirectoryPath(namespace))
}

func (r *fileRepository) DeleteModuleType(namespace string, name string, type_ string) error {
	if err := os.RemoveAll(r.getAbsoluteModuleTypeDirectoryPath(namespace, name, type_)); err != nil {
		return err
	}
	return r.cleanup(r.getAbsoluteModuleNameDirectoryPath(namespace, name))
}

func (r *fileRepository) DeleteModuleVersion(namespace string, name string, type_ string, version string) error {
	filePath := r.getAbsoluteModuleFilePath(namespace, name, type_, version)
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}
	return r.cleanup(r.getAbsoluteModuleTypeDirectoryPath(namespace, name, type_))
}

func (r *fileRepository) cleanup(path string) error {
	splitPath := filepath.SplitList(path)

	for i := len(splitPath) - 1; i <= 0; i-- {
		pathSeg := splitPath[i]

		if pathSeg == modulesDirectory {
			return nil
		}
		subPath := filepath.Join(splitPath[0:i]...)

		if _, err := os.Stat(subPath); os.IsNotExist(err) {
			return nil
		}

		files, err := ioutil.ReadDir(subPath)
		if err != nil {
			return fmt.Errorf("could not list files: %w", err)
		}

		if len(files) == 0 {
			return os.Remove(subPath)
		}

		return nil
	}

	return nil
}

func (r *fileRepository) GetModule(namespace string, name string, type_ string, version string) (module *spec.Module, rerr error) {
	targetAbsModuleFilePath := r.getAbsoluteModuleFilePath(namespace, name, type_, version)

	if _, err := os.Stat(targetAbsModuleFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("not found")
	}

	l := r.newFileLock(targetAbsModuleFilePath)
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	locked, err := l.TryRLockContext(lockCtx, 500*time.Millisecond)
	if !locked || err != nil {
		return nil, fmt.Errorf("could not lock: %s", l.Path())
	}

	defer func() {
		if err := l.Unlock(); err != nil {
			if rerr != nil {
				rerr = fmt.Errorf("%s ; could not unlock: %w", rerr.Error(), err)
			}
			rerr = fmt.Errorf("could not unlock: %w", err)
		}
	}()

	serializedModule, err := ioutil.ReadFile(targetAbsModuleFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read module file: %w", err)
	}

	m := &spec.Module{}
	if err := proto.Unmarshal(serializedModule, m); err != nil {
		return nil, fmt.Errorf("could not unmarhsal proto: %w", err)
	}

	return m, nil
}

func (r *fileRepository) ListModuleNamespaces() ([]string, error) {
	var namespaces []string

	if _, err := os.Stat(r.path); err == nil {
		files, err := ioutil.ReadDir(r.path)
		if err != nil {
			return nil, fmt.Errorf("could not list directories: %w", err)
		}

		for _, f := range files {
			if f.IsDir() {
				namespaces = append(namespaces, f.Name())
			}
		}
	}

	return namespaces, nil
}

func (r *fileRepository) ListModuleNames(namespace string) ([]string, error) {
	var names []string

	directoryPath := r.getAbsoluteModuleNamespaceDirectoryPath(namespace)
	if _, err := os.Stat(directoryPath); err == nil {
		files, err := ioutil.ReadDir(directoryPath)
		if err != nil {
			return nil, fmt.Errorf("could not list directories: %w", err)
		}

		for _, f := range files {
			if f.IsDir() {
				names = append(names, f.Name())
			}
		}
	}

	return names, nil
}

func (r *fileRepository) ListModuleTypes(namespace string, name string) ([]string, error) {
	var types []string

	directoryPath := r.getAbsoluteModuleNameDirectoryPath(namespace, name)
	if _, err := os.Stat(directoryPath); err == nil {
		files, err := ioutil.ReadDir(directoryPath)
		if err != nil {
			return nil, fmt.Errorf("could not list directories: %w", err)
		}

		for _, f := range files {
			if f.IsDir() {
				types = append(types, f.Name())
			}
		}
	}

	return types, nil
}

func (r *fileRepository) ListModuleVersions(namespace string, name string, type_ string) ([]string, error) {
	var versions []string

	directoryPath := r.getAbsoluteModuleTypeDirectoryPath(namespace, name, type_)
	if _, err := os.Stat(directoryPath); err == nil {
		files, err := ioutil.ReadDir(directoryPath)
		if err != nil {
			return nil, fmt.Errorf("could not list directories: %w", err)
		}

		for _, f := range files {
			if strings.HasSuffix(f.Name(), "."+moduleFileExtension) {
				versions = append(versions, strings.TrimSuffix(f.Name(), "."+moduleFileExtension))
			}
		}
	}

	return versions, nil
}
