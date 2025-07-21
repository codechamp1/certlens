package ui

import (
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	docStyle      lipgloss.Style
	errorModal    lipgloss.Style
	sectionHeader lipgloss.Style
	key           lipgloss.Style
	value         lipgloss.Style
}

var Default = Theme{
	docStyle: lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1),

	errorModal: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("#ff5555")),

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

}

func (t Theme) ErrorModalWithWidth(width int) lipgloss.Style {
	return t.errorModal.Width(width)
}
