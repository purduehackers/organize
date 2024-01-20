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

func IntroDescriptionView(width int, text string) string {
	return lipgloss.NewStyle().
		Width(int(math.Round(float64(width)*0.6))).
		Padding(0, 1).
		Render(text) + "\n\n"
}

func PositionListItemView(maxWidth int, title string, description string, selected bool) string {
	titleTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)
	containerStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("63")).
		Width(int(math.Round(float64(maxWidth) * 0.6)))
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

func HomeOptions(viewportWidth int, cursor int) string {
	var maxWidth int = viewportWidth / 2
	events := PositionListItemView(maxWidth, "Events", "View upcoming & past events", cursor == 0)
	organizerPositions := PositionListItemView(maxWidth, "Organize", "View open organizer positions", cursor == 1)

	return lipgloss.JoinHorizontal(lipgloss.Left, events, organizerPositions)
}

func OpenPositionsRows(viewportWidth int, fileNames []string, fileDescriptions []string, cursor int) string {
	var rows []string
	var maxWidth = viewportWidth

	readmeSelected := cursor == 0
	styledReadme := PositionListItemView(maxWidth, fileNames[0], fileDescriptions[0], readmeSelected) + "\n\n\n"
	openPositions := TextWithBackgroundView("#C48FDC", "OPEN POSITIONS", false)
	startHere := styledReadme + openPositions
	rows = append(rows, startHere)

	for i := 1; i < len(fileNames); i++ {
		var row string
		selected := cursor == i
		styledFileName := PositionListItemView(maxWidth, fileNames[i], fileDescriptions[i], selected)
		row = lipgloss.JoinHorizontal(lipgloss.Top, row, styledFileName)
		rows = append(rows, row)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
