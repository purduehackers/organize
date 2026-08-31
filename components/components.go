package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Brand colors from purduehackers.com
const (
	ColorPrimary    = "#7d3bff" // purple
	ColorSecondary  = "#ffee00" // yellow
	ColorSecondaryL = "#ffb700" // amber
	ColorAccent     = "#0084ff" // blue
	ColorAccent2    = "#ff00cc" // magenta
	ColorAccent3    = "#efb9ff" // light purple
	ColorBg         = "#fffbf1" // cream
)

func TextWithBackgroundView(backgroundColor string, textColor string, text string, outerPadding bool) string {
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
		Foreground(lipgloss.Color(textColor)).
		Blink(true)

	return outerContainerStyle.Render(innerContainerStyle.Render(textStyle.Render(text))) + "\n"
}

func IntroDescriptionView(width int) string {
	return lipgloss.NewStyle().
		Width(width).
		Padding(0, 1).
		Render("Purdue Hackers is a group of students who help each other build creative technical projects. We're always looking for a few new organizers to join our team.\n\nThe following roles are open as of August 2026. Apply at https://phack.rs/apply.\n\nGet started at the README. Use arrow keys or vim keys to navigate & enter to select.") + "\n\n"
}

func PositionListItemView(maxWidth int, title string, description string, count string, selected bool) string {
	// Width() excludes the border (2 cells); the inner padding is 2 per side.
	const borderWidth, horizontalPadding = 2, 4
	contentWidth := maxWidth - borderWidth - horizontalPadding
	if contentWidth < 1 {
		contentWidth = 1
	}

	titleTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorAccent2)).
		Bold(true)
	containerStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(ColorPrimary)).
		Width(maxWidth - borderWidth)
	badgeStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(ColorSecondaryL)).
		Foreground(lipgloss.Color("#000")).
		Padding(0, 1).
		MarginLeft(1).
		Bold(true)
	if selected {
		containerStyle = containerStyle.
			BorderForeground(lipgloss.Color(ColorSecondary))
	}
	innerContainerStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(2)

	titleContent := titleTextStyle.Render(title)
	if count != "0" {
		badge := badgeStyle.Render(count)
		titleContent += badge
	}
	descriptionTextContent := lipgloss.NewStyle().Width(contentWidth).Render(description)
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
			BorderForeground(lipgloss.Color(ColorSecondary)).
			Padding(0, 1).
			Bold(true)
	}()

	FooterStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return HeaderStyle.Copy().BorderStyle(b)
	}()
)

// ListLayout is the rendered position list plus the line range each item
// occupies within it, so the caller can scroll the selected item into view.
type ListLayout struct {
	Content     string
	ItemTops    []int
	ItemHeights []int
}

func OpenPositionsListView(width int, fileNames []string, fileDescriptions []string, fileOpenPositionCounts []string, cursor int) ListLayout {
	var b strings.Builder
	tops := make([]int, len(fileNames))
	heights := make([]int, len(fileNames))

	writeItem := func(i int, block string) {
		tops[i] = strings.Count(b.String(), "\n")
		heights[i] = lipgloss.Height(block)
		b.WriteString(block)
		b.WriteString("\n")
	}

	b.WriteString(TextWithBackgroundView(ColorSecondary, "#000000", "ORGANIZE PURDUE HACKERS", true))
	b.WriteString(IntroDescriptionView(width))

	writeItem(0, PositionListItemView(width, fileNames[0], fileDescriptions[0], "0", cursor == 0))
	b.WriteString("\n\n")
	b.WriteString(TextWithBackgroundView(ColorPrimary, ColorBg, "OPEN POSITIONS", false))

	for i := 1; i < len(fileNames); i++ {
		writeItem(i, PositionListItemView(width, fileNames[i], fileDescriptions[i], fileOpenPositionCounts[i], cursor == i))
	}

	return ListLayout{Content: b.String(), ItemTops: tops, ItemHeights: heights}
}
