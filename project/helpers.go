package project

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	ProjectConfigFilename = "wb.yaml"
	ContractsDirectory    = "contracts"
	BuildDirectory        = "build"
	BindingsDirectory     = "bindings"
	MigrationsDirectory   = "migrations"
	TestsDirectory        = "tests"
)

func exists(path string) (bool, error) {
	if path == "" {
		return false, nil
	}

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if !os.IsNotExist(err) {
		return false, err
	}

	return false, nil
}

func FindProject() (*Project, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return findProject(wd)
}

func findProject(path string) (*Project, error) {
	if strings.HasSuffix(path, string(filepath.Separator)) {
		return nil, errors.New("Could not find project root")
	}

	e, err := exists(filepath.Join(path, ProjectConfigFilename))
	if err != nil {
		return nil, err
	}
	if e {
		return NewProjectFromPath(path), nil
	}

	return findProject(filepath.Dir(path))
}
