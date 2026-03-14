package utils

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func BuildTitleLabel(text string) fyne.CanvasObject {
	return widget.NewLabelWithStyle(
		text,
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true,},
	)
}
