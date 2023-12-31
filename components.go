package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func joinPurdueHackersView() string {
	outerContainerStyle := lipgloss.NewStyle().
		Padding(1)
	innerContainerStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Background(lipgloss.
			Color("#fcd34d"))
	textStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Blink(true)

	return outerContainerStyle.Render(innerContainerStyle.Render(textStyle.Render("JOIN PURDUE HACKERS"))) + "\n"
}

func introDescriptionView(width int) string {
	return lipgloss.NewStyle().
		Width(int(math.Round(float64(width)*0.6))).
		Padding(0, 1).
		Render("Purdue Hackers is a group of students who help each other build creative technical projects. We're looking for a few new organizers to join our team.\n\nGet started at the README. Use arrow keys or vim keys to navigate & enter to select.") + "\n\n"
}

func positionListItemView(str string, selected bool) string {
	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)
	containerStyle := lipgloss.NewStyle().
		PaddingLeft(3).
		PaddingRight(3).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("63"))
	if selected {
		containerStyle = containerStyle.
			BorderForeground(lipgloss.Color("#fcd34d"))
	}

	textContent := textStyle.Render(str)
	containerContent := containerStyle.Render(textContent)

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

	for i := 0; i < len(m.fileNames); i += 1 {
		var row string
		for j := 0; j < len(m.fileNames[i]); j += 1 {

			selected := m.cursor_v == i && m.cursor_h == j
			styledFileName := positionListItemView(m.fileNames[i][j], selected)
			row = lipgloss.JoinHorizontal(lipgloss.Top, row, styledFileName)
		}
		rows = append(rows, row)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
