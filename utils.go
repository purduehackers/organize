package main

import (
	"os"

	"github.com/charmbracelet/lipgloss"
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

func renderEntry(str string, selected bool) string {
	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	containerStyle := lipgloss.NewStyle().PaddingLeft(2).PaddingRight(2).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("63"))
	if (selected) {
		containerStyle = lipgloss.NewStyle().PaddingLeft(2).PaddingRight(2).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("226"))
	}

	textContent := textStyle.Render(str)
	containerContent := containerStyle.Render(textContent)

	return containerContent
}