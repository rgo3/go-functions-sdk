package parse

import (
	"go/importer"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/google/logger"
)

// Package holds a golangs package name, directory path
// and exported cloud functions
type Package struct {
	Name       string // golang package name
	Path       string // directory path
	ImportPath string // pkg import path

	Functions []types.Object // exported cloud functions
}

type Packages []Package

// GetPackages parses all toplevel folders in ./ and returns all found packages
// that contain cloud functions
func GetPackages() Packages {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		logger.Fatal(err)
	}

	imp := importer.ForCompiler(token.NewFileSet(), "source", nil)

	defaultImportPath := DefaultImportPath()

	functionPkgs := []Package{}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		importPath := path.Join(defaultImportPath, file.Name())
		importerPkg, err := imp.Import(importPath)
		if err != nil {
			logger.Warningf("Could not import %s:\n%v", importPath, err)
			continue
		}

		pkg := Package{
			Path:       file.Name(),
			ImportPath: importPath,
			Functions:  []types.Object{},
		}

		scope := importerPkg.Scope()
		fnNames := scope.Names()
		for _, name := range fnNames {
			obj := scope.Lookup(name)
			if obj.Exported() && isFunctionBuilder(obj.Type().String()) {
				pkg.Name = obj.Pkg().Name()
				pkg.Functions = append(pkg.Functions, obj)
			}
		}

		if len(pkg.Functions) > 0 {
			functionPkgs = append(functionPkgs, pkg)
		}
	}

	return functionPkgs
}

func (pkgs Packages) Functions() []types.Object {
	functions := []types.Object{}
	for _, pkg := range pkgs {
		functions = append(functions, pkg.Functions...)
	}

	return functions
}

func (pkgs Packages) Paths() []string {
	paths := []string{}
	for _, pkg := range pkgs {
		paths = append(paths, pkg.Path)
	}

	return paths
}

func DefaultImportPath() (importPath string) {
	goPath := os.Getenv("GOPATH")
	pwd := os.Getenv("PWD")
	importPath = strings.TrimPrefix(pwd, path.Join(goPath, "src")+"/")
	return
}

func isFunctionBuilder(t string) bool {
	return strings.HasSuffix(t, "FunctionBuilder")
}
