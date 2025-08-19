package ui

import (
	"github.com/charmbracelet/lipgloss"
)

type ThemeProvider interface {
	DocStyle() lipgloss.Style
	ErrorModalWithWidth(width int) lipgloss.Style
	SectionHeader() lipgloss.Style
	Pane(selected bool, width, height int) lipgloss.Style
	Key(width int) lipgloss.Style
	Value(width int) lipgloss.Style
	Help(width int) lipgloss.Style
}

type Theme struct {
	docStyle      lipgloss.Style
	errorModal    lipgloss.Style
	sectionHeader lipgloss.Style
	key           lipgloss.Style
	value         lipgloss.Style
}

var Default = Theme{
	docStyle: lipgloss.NewStyle(),

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
		Align(lipgloss.Right),

	value: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")),
}

func (t Theme) DocStyle() lipgloss.Style {
	return t.docStyle
}

func (t Theme) ErrorModalWithWidth(width int) lipgloss.Style {
	return t.errorModal.Width(width)
}

func (t Theme) Pane(selected bool, width, height int) lipgloss.Style {
	base := lipgloss.NewStyle().
		Width(width).
		Height(height).Border(lipgloss.NormalBorder()).PaddingLeft(1)

	if selected {
		base = base.BorderForeground(lipgloss.Color("#00BFFF"))
	} else {
		base = base.BorderForeground(lipgloss.Color("#666666"))
	}

	return base
}

func (t Theme) SectionHeader() lipgloss.Style {
	return t.sectionHeader
}

func (t Theme) Key(width int) lipgloss.Style {
	return t.key.Width(width)
}

func (t Theme) Value(width int) lipgloss.Style {
	return t.value.Width(width)
}

func (t Theme) Help(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).MarginLeft(1).Width(width)
}
