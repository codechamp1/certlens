package ui

import (
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	docStyle      lipgloss.Style
	sectionHeader lipgloss.Style
	key           lipgloss.Style
	value         lipgloss.Style
	itemTitle     lipgloss.Style
	itemDesc      lipgloss.Style
}

var Default = Theme{
	docStyle: lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1),

	sectionHeader: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00BFFF")).
		Bold(true).
		Underline(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true),

	key: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFA500")).
		Width(20).
		Align(lipgloss.Right),

	value: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		MaxWidth(50).
		PaddingLeft(1), // Same dark gray as main text for values

	itemTitle: lipgloss.NewStyle().
		Background(lipgloss.Color("#007ACC")). // Consistent rich blue for list title background
		Foreground(lipgloss.Color("#FFFFFF")). // White text
		Padding(0, 1).
		Bold(true),

	itemDesc: lipgloss.NewStyle().
		PaddingLeft(2),
}
