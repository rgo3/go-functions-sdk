package build

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/dergoegge/go-functions-sdk/internal/pkg/parse"
)

const (
	// PluginFolder is the folder in which the plugins are being build
	PluginFolder = "./.build"
)

func removeInitFunc(file string, initFuncDecl *ast.FuncDecl) {
	read, _ := ioutil.ReadFile(file)

	_initFuncText := []byte(strings.Replace(
		string(read[initFuncDecl.Pos()-1:initFuncDecl.End()]),
		"init",
		"_init",
		1,
	))

	withOutInit := read[:initFuncDecl.Pos()-1]
	withOutInit = append(withOutInit, _initFuncText...)
	withOutInit = append(withOutInit, read[initFuncDecl.End():]...)

	err := ioutil.WriteFile(file, withOutInit, 0)
	if err != nil {
		log.Fatal(err)
	}
}

// renames all init funcs in the packages to _init so they wont be called for the created plugins
func removeInitFuncs(pkgs []string) {
	for _, pkg := range pkgs {
		parsedPkgs, _ := parser.ParseDir(token.NewFileSet(), path.Join(PluginFolder, pkg), nil, parser.ParseComments)

		for _, parsedPkg := range parsedPkgs {
			for filePath, file := range parsedPkg.Files {
				for _, decl := range file.Decls {
					switch decl.(type) {
					case *ast.FuncDecl:
						removeInitFunc(filePath, decl.(*ast.FuncDecl))
						break
					}
				}
			}
		}
	}
}

// creates plugin for pkg
func createPlugin(pkg parse.Package, errChan chan error) {
	cmd := exec.Command("go", "build", "-buildmode=plugin",
		"-o", "./"+path.Join(PluginFolder, pkg.Name+".so"),
		"./"+path.Join(PluginFolder, pkg.Path))

	// create output buffer for error logging
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf

	// copy env and turn off go modules
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GO111MODULE=off")

	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("Failed to build plugin for package:%s\n%s", pkg.Name, outBuf.String())
	}

	errChan <- err
}

// creates all plugins for pkgs
func createPlugins(pkgs parse.Packages) error {
	if len(pkgs) == 0 {
		return fmt.Errorf("No pkgs provided to plugin creation")
	}

	errChan := make(chan error)
	for _, pkg := range pkgs {
		go createPlugin(pkg, errChan)
	}

	for range pkgs {
		err := <-errChan
		if err != nil {
			return err
		}
	}

	return nil
}

// Plugins creates golang plugins from all toplevel packages in ./.
// Plugins are named after the package name e.g.: "pkg.so".
// One plugin holds all function builder symbols defined in that package.
func Plugins(pkgs parse.Packages) error {
	// create plugin folder
	os.Mkdir(PluginFolder, os.ModePerm)

	pkgPaths := pkgs.Paths()

	// copy packages
	err := copyPackages(pkgs)
	if err != nil {
		return err
	}

	// remove init func from packages
	removeInitFuncs(pkgPaths)

	// create plugins from packages
	return createPlugins(pkgs)
}
