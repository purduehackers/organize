package utils

import (
	"bufio"
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

type PositionMeta struct {
	FileNames              []string
	FileDescriptions       []string
	FileOpenPositionCounts []string
}

func GetPositionMeta(dir string) (*PositionMeta, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, len(files))
	fileDescriptions := make([]string, len(files))
	fileOpenPositionCounts := make([]string, len(files))
	for i := range files {
		fileName := files[i].Name()
		fileNames[i] = fileName

		file, err := os.Open("directory/" + fileName)
		if err != nil {
			return nil, err
		} else {
			defer file.Close()
			scanner := bufio.NewScanner(file)

			// Read first line, which contains the position description
			scanner.Scan()
			fileDescriptions[i] = scanner.Text()

			// Read the second line, which may or may not contain the number of open positions for this role
			scanner.Scan()
			secondLine := scanner.Text()
			if strings.Contains(secondLine, "Count:") {
				fileOpenPositionCounts[i] = strings.Split(secondLine, " ")[1]
			} else {
				fileOpenPositionCounts[i] = "0"
			}
		}
	}
	positionMetas := PositionMeta{
		FileNames:              fileNames,
		FileDescriptions:       fileDescriptions,
		FileOpenPositionCounts: fileOpenPositionCounts,
	}
	return &positionMetas, nil
}

func Typewrite(s ssh.Session, text string, duration time.Duration) {
	for _, char := range text {
		fmt.Fprint(s, string(char))
		time.Sleep(duration * time.Millisecond)
	}
	fmt.Fprint(s, "\n")
}
