package build

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/dergoegge/go-functions-sdk/internal/pkg/parse"
	"github.com/google/logger"
)

const (
	PluginFolder = "./plugins"
)

func findFiles(dirPath string) []string {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	result := make([]string, 0)
	for _, f := range files {
		pkg := path.Join(dirPath, f.Name())
		if f.IsDir() {
			//result = append(result, findFiles(pkg)...)
			continue
		}

		result = append(result, pkg)
	}

	return result
}

func copyFile(src string, dst string) error {
	err := os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	info, err := s.Stat()
	if err != nil {
		return err
	}

	err = os.Chmod(f.Name(), info.Mode())
	if err != nil {
		return err
	}

	_, err = io.Copy(f, s)
	return err
}

func removeInitFunc(file string, initFuncDecl *ast.FuncDecl) {
	logger.Infof("Removing init from %s\n", file)
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

// renames all init funcs in the packages to _init so they wont be called
// in for the created plugins
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

func createPlugins(pkgs parse.Packages) error {
	if len(pkgs) == 0 {
		return fmt.Errorf("No pkgs provided to plugin creation")
	}

	errChan := make(chan error)

	for _, pkg := range pkgs {
		cmd := exec.Command("go", "build", "-buildmode=plugin",
			"-o", "./"+path.Join(PluginFolder, pkg.Name+".so"),
			"./"+path.Join(PluginFolder, pkg.Path))

		logger.Info("Creating plugin for package:", pkg.Name)
		go func(pkg string) {
			errChan <- cmd.Run()
			logger.Info("Successfully created plugin for package:", pkg)
		}(pkg.Name)
	}

	for range pkgs {
		err := <-errChan
		if err != nil {
			logger.Error("Failed to create plugin", err)
			return err
		}
	}

	return nil
}

func copyPackages(pkgs parse.Packages) error {
	for _, pkg := range pkgs {
		files := findFiles(pkg.Path)

		for _, file := range files {
			filePath := path.Join(PluginFolder, file)
			copyFile(file, filePath)

			read, _ := ioutil.ReadFile(filePath)
			newContent := strings.Replace(string(read), "package "+pkg.Name, "package main", 1)
			err := ioutil.WriteFile(filePath, []byte(newContent), 0)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Plugins creates golang plugins from all toplevel packages in ./
func Plugins(pkgs parse.Packages) error {
	// remove previous plugin foulder
	os.RemoveAll(PluginFolder)

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
