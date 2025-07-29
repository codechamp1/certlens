package ui

import "github.com/charmbracelet/lipgloss"

type uiLayout struct {
	TotalWidth     int
	TotalHeight    int
	UsableWidth    int
	UsableHeight   int
	LeftPaneWidth  int
	RightPaneWidth int
}

func calculateLayout(width, height int, docStyle lipgloss.Style) uiLayout {
	hPadding, vPadding := docStyle.GetFrameSize()

	// for border top
	usableH := height - vPadding - 3
	usableW := width - hPadding

	// for border left and right
	left := usableW/2 - 2
	right := usableW - left - 4

	return uiLayout{
		LeftPaneWidth:  left,
		RightPaneWidth: right,
		UsableWidth:    usableW,
		UsableHeight:   usableH,
		TotalHeight:    height,
		TotalWidth:     width,
	}
}
