package build

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/dergoegge/go-functions-sdk/internal/pkg/parse"
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
