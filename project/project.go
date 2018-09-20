package project

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var srcPaths []string

func init() {
	envGoPath := os.Getenv("GOPATH")
	goPaths := filepath.SplitList(envGoPath)
	if len(goPaths) == 0 {
		fmt.Println("$GOPATH is not set")
		os.Exit(1)
	}
	srcPaths = make([]string, 0, len(goPaths))
	for _, goPath := range goPaths {
		srcPaths = append(srcPaths, filepath.Join(goPath, "src"))
	}
}

type Project struct {
	absPath string
	srcPath string
	license License
	name    string
}

func NewProject(projectName string) *Project {
	if projectName == "" {
		return nil
	}

	p := new(Project)
	p.name = projectName

	p.absPath = findPackage(projectName)

	if p.absPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil
		}
		for _, srcPath := range srcPaths {
			goPath := filepath.Dir(srcPath)
			if filepathHasPrefix(wd, goPath) {
				p.absPath = filepath.Join(srcPath, projectName)
				break
			}
		}
	}

	if p.absPath == "" {
		p.absPath = filepath.Join(srcPaths[0], projectName)
	}

	return p
}

func findPackage(packageName string) string {
	if packageName == "" {
		return ""
	}

	for _, srcPath := range srcPaths {
		packagePath := filepath.Join(srcPath, packageName)
		if e, _ := exists(packagePath); e {
			return packagePath
		}
	}

	return ""
}

func NewProjectFromPath(absPath string) *Project {
	if absPath == "" || !filepath.IsAbs(absPath) {
		return nil
	}

	p := new(Project)
	p.absPath = absPath
	p.name = filepath.ToSlash(trimSrcPath(p.absPath, p.SrcPath()))
	return p
}

func trimSrcPath(absPath, srcPath string) string {
	relPath, err := filepath.Rel(srcPath, absPath)
	if err != nil {
		fmt.Println("Cobra supports project only within $GOPATH: " + err.Error())
		os.Exit(1)
	}
	return relPath
}

func (p *Project) License() License {
	if p.license.Text == "" && p.license.Name != "None" {
		p.license = getLicense()
	}

	return p.license
}

func (p Project) Name() string {
	return p.name
}

func (p Project) AbsPath() string {
	return p.absPath
}

func (p *Project) SrcPath() string {
	if p.srcPath != "" {
		return p.srcPath
	}

	if p.absPath == "" {
		p.srcPath = srcPaths[0]
		return p.srcPath
	}

	for _, srcPath := range srcPaths {
		if filepathHasPrefix(p.absPath, srcPath) {
			p.srcPath = srcPath
			break
		}
	}

	return p.srcPath
}

func (p *Project) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"project":   p.Name(),
		"license":   p.License(),
		"copyright": copyrightLine(),
	}
}

func filepathHasPrefix(path string, prefix string) bool {
	if len(path) <= len(prefix) {
		return false
	}

	if runtime.GOOS == "windows" {
		return strings.EqualFold(path[0:len(prefix)], prefix)
	}

	return path[0:len(prefix)] == prefix
}
