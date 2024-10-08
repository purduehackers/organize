package components

import (
	"math"

	"github.com/charmbracelet/lipgloss"
)

func TextWithBackgroundView(backgroundColor string, text string, outerPadding bool) string {
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

func IntroDescriptionView(width int) string {
	return lipgloss.NewStyle().
		Width(int(math.Round(float64(width)*0.6))).
		Padding(0, 1).
		Render("Purdue Hackers is a group of students who help each other build creative technical projects. We're looking for a few new organizers to join our team during the Fall 2024 semester.\n\nThe following positions are open as of August 2024.\n\nGet started at the README. Use arrow keys or vim keys to navigate & enter to select.") + "\n\n"
}

func PositionListItemView(maxWidth int, title string, description string, count string, selected bool) string {
	titleTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)
	containerStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("63")).
		Width(int(math.Round(float64(maxWidth) * 0.6)))
	badgeStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#fcd34d")).
		Foreground(lipgloss.Color("#000")).
		Padding(0, 1).
		MarginLeft(1).
		Bold(true)
	if selected {
		containerStyle = containerStyle.
			BorderForeground(lipgloss.Color("#fcd34d"))
	}
	innerContainerStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(2)

	titleContent := titleTextStyle.Render(title)
	if count != "0" {
		badge := badgeStyle.Render(count)
		titleContent += badge
	}
	descriptionTextContent := lipgloss.NewStyle().Render(description)
	textContent := titleContent + "\n" + descriptionTextContent

	innerContainerContent := innerContainerStyle.Render(textContent)
	containerContent := containerStyle.Render(innerContainerContent)

	return containerContent
}

var (
	HeaderStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().
			BorderStyle(b).
			BorderForeground(lipgloss.Color("#fcd34d")).
			Padding(0, 1).
			Bold(true)
	}()

	FooterStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return HeaderStyle.Copy().BorderStyle(b)
	}()
)

func OpenPositionsGrid(width int, fileNames []string, fileDescriptions []string, fileOpenPositionCounts []string, cursor int) string {
	var rows []string
	var maxWidth = width

	readmeSelected := cursor == 0
	styledReadme := PositionListItemView(maxWidth, fileNames[0], fileDescriptions[0], "0", readmeSelected) + "\n\n\n"
	openPositions := TextWithBackgroundView("#C48FDC", "OPEN POSITIONS", false)
	startHere := styledReadme + openPositions
	rows = append(rows, startHere)

	for i := 1; i < len(fileNames); i++ {
		var row string
		selected := cursor == i
		styledFileName := PositionListItemView(maxWidth, fileNames[i], fileDescriptions[i], fileOpenPositionCounts[i], selected)
		row = lipgloss.JoinHorizontal(lipgloss.Top, row, styledFileName)
		rows = append(rows, row)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
