package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/ssh"
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

func typewrite(s ssh.Session, text string, duration time.Duration) {
	for _, char := range text {
		fmt.Fprint(s, string(char))
		time.Sleep(duration * time.Millisecond)
	}
	fmt.Fprint(s, "\n")
}
