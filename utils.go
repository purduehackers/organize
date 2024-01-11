package main

import (
	"bufio"
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

type positionMeta struct {
	fileNames        []string
	fileDescriptions []string
}

func getPositionMeta(dir string) (*positionMeta, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, len(files))
	fileDescriptions := make([]string, len(files))
	for i := range files {
		fileName := files[i].Name()
		fileNames[i] = fileName

		file, err := os.Open("directory/" + fileName)
		if err != nil {
			return nil, err
		} else {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			scanner.Scan()
			fileDescriptions[i] = scanner.Text()
		}
	}
	positionMetas := positionMeta{
		fileNames:        fileNames,
		fileDescriptions: fileDescriptions,
	}
	return &positionMetas, nil
}

func typewrite(s ssh.Session, text string, duration time.Duration) {
	for _, char := range text {
		fmt.Fprint(s, string(char))
		time.Sleep(duration * time.Millisecond)
	}
	fmt.Fprint(s, "\n")
}
