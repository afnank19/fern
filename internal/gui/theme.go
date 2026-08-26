package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Fern's dark palette. Tune these values to retune the whole UI; anything not
// overridden here falls through to Fyne's built-in dark variant.
//
// NOTE: the embedded DefaultTheme is always queried with VariantDark so the
// app stays dark even if the OS prefers light.
var (
	colBackground      = color.RGBA{R: 0x11, G: 0x13, B: 0x18, A: 0xff}
	colPanel           = color.RGBA{R: 0x17, G: 0x1a, B: 0x20, A: 0xff} // menus, tab bar
	colInputBackground = color.RGBA{R: 0x1e, G: 0x21, B: 0x28, A: 0xff}
	colButton          = color.RGBA{R: 0x26, G: 0x2a, B: 0x33, A: 0xff}
	colHover           = color.RGBA{R: 0x2f, G: 0x34, B: 0x40, A: 0xff}
	colPressed         = color.RGBA{R: 0x1b, G: 0x1e, B: 0x24, A: 0xff}
	colDisabled        = color.RGBA{R: 0x56, G: 0x5b, B: 0x64, A: 0xff}
	colDisabledButton  = color.RGBA{R: 0x1d, G: 0x20, B: 0x26, A: 0xff}
	colForeground      = color.RGBA{R: 0xe6, G: 0xe6, B: 0xe6, A: 0xff}
	colPlaceholder     = color.RGBA{R: 0x6b, G: 0x70, B: 0x78, A: 0xff}
	colSeparator       = color.RGBA{R: 0x26, G: 0x2a, B: 0x33, A: 0xff}
	colInputBorder     = color.RGBA{R: 0x33, G: 0x39, B: 0x45, A: 0xff}
	colScrollBar       = color.RGBA{R: 0x3a, G: 0x3f, B: 0x4a, A: 0xff}
	colSelection       = color.RGBA{R: 0x1f, G: 0x3a, B: 0x29, A: 0xff}
	colShadow          = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x40}
	colOverlay         = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x99}

	// Accent: fern green. Bright enough that on-primary text must be dark.
	colAccent             = color.RGBA{R: 0x4a, G: 0xde, B: 0x80, A: 0xff}
	colForegroundOnAccent = color.RGBA{R: 0x07, G: 0x13, B: 0x0c, A: 0xff}
)

// fernTheme implements fyne.Theme on top of the default theme: colors come
// from the dark palette above, every font is forced to the bundled monospace,
// and corner radii are flattened to square. Everything else delegates.
type fernTheme struct {
	fyne.Theme
}

func newFernTheme() fyne.Theme {
	return &fernTheme{Theme: theme.DefaultTheme()}
}

func (t *fernTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return colBackground
	case theme.ColorNameMenuBackground, theme.ColorNameHeaderBackground:
		return colPanel
	case theme.ColorNameOverlayBackground:
		return colOverlay
	case theme.ColorNameInputBackground:
		return colInputBackground
	case theme.ColorNameInputBorder:
		return colInputBorder
	case theme.ColorNameButton:
		return colButton
	case theme.ColorNameHover:
		return colHover
	case theme.ColorNamePressed:
		return colPressed
	case theme.ColorNameDisabled:
		return colDisabled
	case theme.ColorNameDisabledButton:
		return colDisabledButton
	case theme.ColorNameForeground:
		return colForeground
	case theme.ColorNamePlaceHolder:
		return colPlaceholder
	case theme.ColorNameSeparator:
		return colSeparator
	case theme.ColorNameScrollBar:
		return colScrollBar
	case theme.ColorNameScrollBarBackground:
		return colBackground
	case theme.ColorNameShadow:
		return colShadow
	case theme.ColorNameSelection:
		return colSelection
	case theme.ColorNamePrimary, theme.ColorNameFocus:
		return colAccent
	case theme.ColorNameHyperlink:
		return color.RGBA{R: 0x58, G: 0xa6, B: 0xff, A: 0xff}
	case theme.ColorNameError:
		return color.RGBA{R: 0xef, G: 0x53, B: 0x50, A: 0xff}
	case theme.ColorNameWarning:
		return color.RGBA{R: 0xff, G: 0xb7, B: 0x4d, A: 0xff}
	case theme.ColorNameSuccess:
		return color.RGBA{R: 0x66, G: 0xbb, B: 0x6a, A: 0xff}
	case theme.ColorNameForegroundOnPrimary:
		return colForegroundOnAccent
	}
	return t.Theme.Color(name, theme.VariantDark)
}

// Font serves the bundled Noto Sans Mono for every style while no custom
// fonts are installed. To switch later: embed real TTFs and dispatch on the
// style fields (Bold, Italic, Monospace...) instead of this blanket override.
func (t *fernTheme) Font(_ fyne.TextStyle) fyne.Resource {
	return t.Theme.Font(fyne.TextStyle{Monospace: true})
}

// Size flattens every rounded corner (buttons, inputs, selection highlights,
// scrollbars); other sizes stay at Fyne defaults.
func (t *fernTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameInputRadius, theme.SizeNameSelectionRadius, theme.SizeNameScrollBarRadius:
		return 0
	}
	return t.Theme.Size(name)
}
