package main

import (
	"github.com/charmbracelet/lipgloss"
)

func JoinPurdueHackers() string {
	outerContainerStyle := lipgloss.NewStyle().Padding(1)
	innerContainerStyle := lipgloss.NewStyle().Padding(0, 1).Background(lipgloss.Color("#fcd34d"))
	textStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#000000")).Blink(true)

	return outerContainerStyle.Render(innerContainerStyle.Render(textStyle.Render("JOIN PURDUE HACKERS"))) + "\n"
}

func IntroDescription(width int) string {
	return lipgloss.NewStyle().Width(width).Padding(0, 1).Render("Purdue Hackers is a group of students who help each other build creative technical projects. We're looking for a few new organizers to join our team.\n\nGet started at the README.") + "\n\n"
}

func PositionListItem(str string, selected bool) string {
	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)
	containerStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))
	if selected {
		containerStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingRight(2).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("226"))
	}

	textContent := textStyle.Render(str)
	containerContent := containerStyle.Render(textContent)

	return containerContent
}
