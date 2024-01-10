package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func textWithBackgroundView(backgroundColor string, text string, outerPadding bool) string {
	outerContainerStyle := lipgloss.NewStyle()
	if outerPadding {
		outerContainerStyle = outerContainerStyle.Padding(1)
	}
	innerContainerStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Background(lipgloss.
			Color(backgroundColor))
	textStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Blink(true)

	return outerContainerStyle.Render(innerContainerStyle.Render(textStyle.Render(text))) + "\n"
}

func introDescriptionView(width int) string {
	return lipgloss.NewStyle().
		Width(int(math.Round(float64(width)*0.6))).
		Padding(0, 1).
		Render("Purdue Hackers is a group of students who help each other build creative technical projects. We're looking for a few new organizers to join our team.\n\nGet started at the README. Use arrow keys or vim keys to navigate & enter to select.") + "\n\n"
}

func positionListItemView(maxWidth int, title string, description string, selected bool) string {
	titleTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)
	containerStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("63")).
		Width(int(math.Round(float64(maxWidth)*0.6)))
	if selected {
		containerStyle = containerStyle.
			BorderForeground(lipgloss.Color("#fcd34d"))
	}
	innerContainerStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(2)

	titleContent := titleTextStyle.Render(title)
	descriptionTextContent := lipgloss.NewStyle().Render(description)
	textContent := titleContent + "\n" + descriptionTextContent

	innerContainerContent := innerContainerStyle.Render(textContent)
	containerContent := containerStyle.Render(innerContainerContent)

	return containerContent
}

var (
	headerStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().
			BorderStyle(b).
			BorderForeground(lipgloss.Color("#fcd34d")).
			Padding(0, 1).
			Bold(true)
	}()

	footerStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return headerStyle.Copy().BorderStyle(b)
	}()
)

func (m model) headerView() string {
	title := headerStyle.Render(m.selectedFileName)
	line := strings.Repeat(lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fcd34d")).
		Render("─"), Max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	helpView := lipgloss.PlaceHorizontal(m.viewport.Width, lipgloss.Right, m.help.View(m.keys))

	info := footerStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat(lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fcd34d")).
		Render("─"), Max(0, m.viewport.Width-lipgloss.Width(info)))
	footerInfo := lipgloss.JoinHorizontal(lipgloss.Center, line, info)

	return helpView + "\n" + footerInfo
}

func (m model) openPositionsGrid() string {
	var rows []string
	var maxWidth = m.viewport.Width

	readmeSelected := m.cursor == 0
	styledReadme := positionListItemView(maxWidth, m.fileNames[0], m.fileDescriptions[0], readmeSelected) + "\n\n\n"
	openPositions := textWithBackgroundView("#C48FDC", "OPEN POSITIONS", false)
	startHere := styledReadme + openPositions
	rows = append(rows, startHere)

	for i := 1; i < len(m.fileNames); i++ {
		var row string
		selected := m.cursor == i
		styledFileName := positionListItemView(maxWidth, m.fileNames[i], m.fileDescriptions[i], selected)
		row = lipgloss.JoinHorizontal(lipgloss.Top, row, styledFileName)
		rows = append(rows, row)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
