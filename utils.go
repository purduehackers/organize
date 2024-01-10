package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/ssh"
)

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func getFileNames(dir string) ([]string, error) {
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

func readFirstLines(dir string) ([]string, error) {
	fileNames, _ := getFileNames(dir)
	firstLines := make([]string, len(fileNames))

	for i := 0; i < len(fileNames); i++ {
		// file, err := os.Open("directory/" + fileNames[i])
		// if err != nil {
		// 	return nil, err
		// } else {
		// 	defer file.Close()
		// 	scanner := bufio.NewScanner(file)
		// 	firstLines[i] = scanner.Text()
		// }

		content, err := os.ReadFile("directory/" + fileNames[i])
		if err != nil {
			return nil, err
		} else {
			fileContent := string(content)
			firstLines[i] = strings.Split(fileContent, "\n")[0]
		}
	}

	return firstLines, nil
}

func typewrite(s ssh.Session, text string, duration time.Duration) {
	for _, char := range text {
		fmt.Fprint(s, string(char))
		time.Sleep(duration * time.Millisecond)
	}
	fmt.Fprint(s, "\n")
}
