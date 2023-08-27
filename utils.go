package main

import (
	"os"
)

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func readFiles(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, len(files))
	for i := range files {
		fileNames[i] = files[i].Name()
	}

	return fileNames, nil
}
