package ui

type modelLayout struct {
	TotalWidth     int
	TotalHeight    int
	UsableWidth    int
	UsableHeight   int
	LeftPaneWidth  int
	RightPaneWidth int
}

func (m *Model) updateLayout(width, height int) {
	hPadding, vPadding := m.theme.DocStyle().GetFrameSize()

	// for border top
	usableH := height - vPadding - 2
	usableW := width - hPadding

	// for border left and right
	left := usableW/2 - 2
	right := usableW - left - 4

	m.layout = modelLayout{
		LeftPaneWidth:  left,
		RightPaneWidth: right,
		UsableWidth:    usableW,
		UsableHeight:   usableH,
		TotalHeight:    height,
		TotalWidth:     width,
	}

	m.secrets.SetSize(left, usableH)
	m.inspectedViewport.Width = right
	m.inspectedViewport.Height = usableH
}
