package ui

import (
	"strings"
)

type HelpViewModel struct {
	pane  Pane
	theme ThemeProvider
	width int
}

func NewHelpViewModel(p Pane, tp ThemeProvider) HelpViewModel {
	return HelpViewModel{
		pane:  p,
		theme: tp,
	}
}

type keyHint struct {
	key, desc string
}

var baseKeyHints = []keyHint{
	{"u", "refresh"},
	{"tab", "switch pane"},
	{"p", "switch pane"},
	{"r", "toggle raw"},
	{"c", "copy cert"},
	{"C", "copy key"},
	{"q", "quit"},
}

var leftPaneKeyHints = []keyHint{
	{"↑/↓", "navigate"},
	{"←/→", "switch list page"},
	{"/", "filter"},
}

var rightPaneKeyHints = []keyHint{
	{"↑/↓", "scroll"},
	{"←/→", "switch cert page"},
	{"enter", "select"},
}

const separator = "  •  "

func (h HelpViewModel) View() string {
	var keyHints []keyHint
	switch h.pane {
	case LeftPane:
		keyHints = append(leftPaneKeyHints, baseKeyHints...)
	case RightPane:
		keyHints = append(rightPaneKeyHints, baseKeyHints...)
	}

	return h.theme.Help(h.width).Render(formatKeyHints(keyHints))
}

func formatKeyHints(hints []keyHint) string {
	var formattedHints []string
	for _, hint := range hints {
		formattedHints = append(formattedHints, hint.key+": "+hint.desc)
	}
	return strings.Join(formattedHints, separator)
}

func (h *HelpViewModel) SetPane(p Pane) {
	h.pane = p
}

func (h *HelpViewModel) SetWidth(width int) {
	h.width = width
}
