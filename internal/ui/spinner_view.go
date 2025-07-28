package ui

import "github.com/charmbracelet/bubbles/spinner"

func NewSpinnerViewModel() SpinnerViewModel {
	return SpinnerViewModel{
		spinner: spinner.New(),
	}
}
