package core

import (
	"fmt"
	"go/build"
	"path/filepath"
	"strings"
)

// PackageInfo retrives go package info based on path
type PackageInfo struct {
	Path   string
	Short  string
	Full   string
	Deps   string
	Cmd    string
	Build  string
	Output string
}

const defaultBuild = "__build"

// GetPackageInfo returns the package info based on path
func GetPackageInfo(path string, cmd string, deps string) (PackageInfo, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return PackageInfo{}, err
	}

	goSrc := filepath.Join(build.Default.GOPATH, "src")
	rel, err := filepath.Rel(goSrc, abs)

	if err != nil || strings.HasPrefix(rel, ".") {
		return PackageInfo{}, fmt.Errorf("%v not go package", abs)
	}

	cmdFull := filepath.Join(rel, cmd)

	return PackageInfo{
		Path:   abs,
		Full:   rel,
		Short:  filepath.Base(rel),
		Deps:   strings.Replace(deps, ",", " ", -1),
		Cmd:    cmdFull,
		Build:  defaultBuild,
		Output: filepath.Base(cmdFull),
	}, nil
}
