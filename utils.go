package main

import (
	"os"
	"path/filepath"
)

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func readFiles(dir string) ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(filepath.Join(wd, dir))
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, len(files))
	for i := range files {
		fileNames[i] = files[i].Name()
	}

	return fileNames, nil
}
