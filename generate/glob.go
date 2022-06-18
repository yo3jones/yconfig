package generate

import (
	"os"
	"path/filepath"
)

func isDir(file string) bool {
	info, err1 := os.Stat(file)
	if err1 != nil {
		panic(err1)
	}
	return info.IsDir()
}

func findFiles(glob string) []string {
	found, err1 := filepath.Glob(glob)
	if err1 != nil {
		panic(err1)
	}

	files := []string{}
	for _, f := range found {
		if !isDir(f) {
			files = append(files, f)
		}
	}
	return files
}

func glob(root string, include, exclude []string) []string {
	filesSet := map[string]bool{}

	for _, inc := range include {
		files := findFiles(filepath.Join(root, inc))
		for _, file := range files {
			filesSet[file] = true
		}
	}

	for _, excl := range exclude {
		files := findFiles(filepath.Join(root, excl))
		for _, file := range files {
			delete(filesSet, file)
		}
	}

	files := []string{}
	for file := range filesSet {
		files = append(files, file)
	}

	return files
}
