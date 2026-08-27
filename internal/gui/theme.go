package gui

import (
	_ "embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed assets/fonts/IoskeleyMono-Regular.ttf
var monoFontRegular []byte

//go:embed assets/fonts/IoskeleyMono-Bold.ttf
var monoFontBold []byte

//go:embed assets/fonts/IoskeleyMono-Light.ttf
var monoFontLight []byte

//go:embed assets/fonts/IoskeleyMono-Medium.ttf
var monoFontMedium []byte

var (
	fontRegular = &fyne.StaticResource{
		StaticName:    "IoskeleyMono-Regular.ttf",
		StaticContent: monoFontRegular,
	}
	fontBold = &fyne.StaticResource{
		StaticName:    "IoskeleyMono-Bold.ttf",
		StaticContent: monoFontBold,
	}
	fontLight = &fyne.StaticResource{
		StaticName:    "IoskeleyMono-Light.ttf",
		StaticContent: monoFontLight,
	}
	fontMedium = &fyne.StaticResource{
		StaticName:    "IoskeleyMono-Medium.ttf",
		StaticContent: monoFontMedium,
	}
)

// Fern's editor-style dark palette.
//
// The reference UI is much closer to a native hex/debugger tool than a
// conventional "dark application" theme:
//
//   - almost-black workspaces
//   - #232323-ish panels
//   - restrained steel-blue accent
//   - neutral gray controls
//   - strong but subtle separators
//   - square geometry
//
// Keep these centralized so the entire application can be retuned without
// touching individual widgets.
var (
	// -------------------------------------------------------------------------
	// Surfaces
	// -------------------------------------------------------------------------

	// Main application/workspace background.
	colBackground = color.RGBA{
		R: 0x0f, G: 0x0f, B: 0x0f, A: 0xff,
	}

	// Menus, headers, tab strips, toolbars and secondary panels.
	colPanel = color.RGBA{
		R: 0x23, G: 0x23, B: 0x23, A: 0xff,
	}

	// Text fields, inspectors and other inset controls.
	colInputBackground = color.RGBA{
		R: 0x1a, G: 0x1a, B: 0x1a, A: 0xff,
	}

	// Normal buttons / inactive controls.
	colButton = color.RGBA{
		R: 0x2b, G: 0x2b, B: 0x2b, A: 0xff,
	}

	// Mouse-over state.
	colHover = color.RGBA{
		R: 0x35, G: 0x35, B: 0x35, A: 0xff,
	}

	// Pressed state.
	colPressed = color.RGBA{
		R: 0x1e, G: 0x1e, B: 0x1e, A: 0xff,
	}

	// Disabled text/icon state.
	colDisabled = color.RGBA{
		R: 0x66, G: 0x66, B: 0x66, A: 0xff,
	}

	// Disabled button surface.
	colDisabledButton = color.RGBA{
		R: 0x1d, G: 0x1d, B: 0x1d, A: 0xff,
	}

	// -------------------------------------------------------------------------
	// Text
	// -------------------------------------------------------------------------

	colForeground = color.RGBA{
		R: 0xd8, G: 0xd8, B: 0xd8, A: 0xff,
	}

	colPlaceholder = color.RGBA{
		R: 0x75, G: 0x75, B: 0x75, A: 0xff,
	}

	// -------------------------------------------------------------------------
	// Borders / separators
	// -------------------------------------------------------------------------

	colSeparator = color.RGBA{
		R: 0x3d, G: 0x3d, B: 0x3d, A: 0xff,
	}

	colInputBorder = color.RGBA{
		R: 0x3a, G: 0x3a, B: 0x3a, A: 0xff,
	}

	// -------------------------------------------------------------------------
	// Scrolling / overlays
	// -------------------------------------------------------------------------

	colScrollBar = color.RGBA{
		R: 0x4a, G: 0x4a, B: 0x4a, A: 0xff,
	}

	// Scrollbar track.
	colScrollBarBackground = color.RGBA{
		R: 0x12, G: 0x12, B: 0x12, A: 0xff,
	}

	// Selection in editors/tables.
	//
	// Muted blue rather than green. This is intentionally darker than the
	// primary accent so selected text does not become visually overwhelming.
	colSelection = color.RGBA{
		R: 0x2c, G: 0x42, B: 0x5e, A: 0xff,
	}

	colShadow = color.RGBA{
		R: 0x00, G: 0x00, B: 0x00, A: 0x60,
	}

	colOverlay = color.RGBA{
		R: 0x08, G: 0x08, B: 0x08, A: 0xe0,
	}

	// -------------------------------------------------------------------------
	// Accent
	// -------------------------------------------------------------------------

	// Steel blue from the reference UI.
	//
	// This is intentionally much less saturated than the old fern-green
	// accent. It works better for tabs, focus indicators and primary controls.
	colAccent = color.RGBA{
		R: 0x3f, G: 0x6f, B: 0xa3, A: 0xff,
	}

	// The reference uses light text on selected/active blue controls.
	colForegroundOnAccent = color.RGBA{
		R: 0xf0, G: 0xf2, B: 0xf5, A: 0xff,
	}

	// -------------------------------------------------------------------------
	// Semantic colors
	// -------------------------------------------------------------------------

	colHyperlink = color.RGBA{
		R: 0x58, G: 0xa6, B: 0xff, A: 0xff,
	}

	colError = color.RGBA{
		R: 0xe5, G: 0x5c, B: 0x5c, A: 0xff,
	}

	colWarning = color.RGBA{
		R: 0xd6, G: 0xa8, B: 0x4f, A: 0xff,
	}

	colSuccess = color.RGBA{
		R: 0x65, G: 0xb3, B: 0x78, A: 0xff,
	}
)

// fernTheme implements fyne.Theme on top of Fyne's default theme.
//
// Only the visual pieces that matter for the editor aesthetic are overridden.
// Everything else continues to come from Fyne, which keeps this theme
// relatively resistant to Fyne version changes.
type fernTheme struct {
	fyne.Theme
}

func newFernTheme() fyne.Theme {
	return &fernTheme{
		Theme: theme.DefaultTheme(),
	}
}

func (t *fernTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return colBackground

	case theme.ColorNameMenuBackground,
		theme.ColorNameHeaderBackground:
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
		return colScrollBarBackground

	case theme.ColorNameShadow:
		return colShadow

	case theme.ColorNameSelection:
		return colSelection

	case theme.ColorNamePrimary,
		theme.ColorNameFocus:
		return colAccent

	case theme.ColorNameHyperlink:
		return colHyperlink

	case theme.ColorNameError:
		return colError

	case theme.ColorNameWarning:
		return colWarning

	case theme.ColorNameSuccess:
		return colSuccess

	case theme.ColorNameForegroundOnPrimary:
		return colForegroundOnAccent
	}

	// Always fall back to the dark Fyne variant. This is important because the
	// application should remain dark even if the operating system is using a
	// light theme.
	return t.Theme.Color(name, theme.VariantDark)
}

// Return the embedded IoskeleyMono font, picking weight by style.
func (t *fernTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Bold {
		return fontBold
	}
	return fontRegular
}

// Flatten Fyne's rounded UI geometry.
//
// The reference application is visually much closer to a traditional desktop
// tool: inputs, selections and scrollbars have essentially square corners.
func (t *fernTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameInputRadius,
		theme.SizeNameSelectionRadius,
		theme.SizeNameScrollBarRadius:
		return 0

	case theme.SizeNameInnerPadding:
		return 4
	}

	return t.Theme.Size(name)
}
