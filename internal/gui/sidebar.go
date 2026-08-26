package gui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// SidebarCallbacks decouple the sidebar from image state: it reports user
// intent, and the owner decides what happens.
type SidebarCallbacks struct {
	OnParam  func(kind OpKind, key string, val float64) // slider/check moved
	OnCommit func(kind OpKind)                          // an op's Apply button pressed
	OnUndo   func()                                     // undo last committed op
}

// Sidebar renders editing controls generated from the op registry. It holds
// no image state: values flow out through callbacks, and ResetOp/Reset let
// the owner snap controls back to defaults (after commit or image load).
type Sidebar struct {
	container *fyne.Container
	cb        SidebarCallbacks

	controls map[ctlKey]ctlWidgets
}

type ctlKey struct {
	kind OpKind
	key  string
}

type ctlWidgets struct {
	label  *widget.Label
	slider *widget.Slider
	check  *widget.Check
}

func NewSidebar(cb SidebarCallbacks) *Sidebar {
	s := &Sidebar{
		cb:       cb,
		controls: make(map[ctlKey]ctlWidgets),
	}
	s.container = container.NewVBox(
		widget.NewButtonWithIcon("Undo", theme.ContentUndoIcon(), s.cb.OnUndo),
		s.buildTabs(),
	)
	return s
}

// buildTabs generates tabs from the registry. Tab order follows first
// appearance in registryOrder; control order follows each op's declared
// Params. Tabs hold plain VBoxes: scrolling is handled once, by the outer
// scroll in CanvasObject.
func (s *Sidebar) buildTabs() fyne.CanvasObject {
	tabOrder := make([]string, 0, 2)
	catContents := make(map[string][]fyne.CanvasObject)

	for _, kind := range registryOrder {
		def := registry[kind]

		if _, seen := catContents[def.Category]; !seen {
			tabOrder = append(tabOrder, def.Category)
		}

		catContents[def.Category] = append(
			catContents[def.Category],
			s.buildOpControls(kind, def)...,
		)
	}

	tabs := container.NewAppTabs()

	for _, cat := range tabOrder {
		tabs.Append(
			container.NewTabItem(
				cat,
				container.NewVBox(catContents[cat]...),
			),
		)
	}

	tabs.SetTabLocation(container.TabLocationTop)

	return tabs
}

// buildOpControls builds one visually separated section for an operation.
func (s *Sidebar) buildOpControls(kind OpKind, def opDef) []fyne.CanvasObject {
	// Section contents.
	items := make([]fyne.CanvasObject, 0, len(def.Params)*2+4)

	// Bold section title.
	title := widget.NewLabel(def.Label)
	title.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	items = append(
		items,
		title,
	)

	for _, p := range def.Params {
		k := ctlKey{kind: kind, key: p.Key}

		label := widget.NewLabel(formatParam(p, p.Default))
		w := ctlWidgets{
			label: label,
		}

		switch p.Widget {
		case CheckWidget:
			check := widget.NewCheck(p.Label, nil)
			check.SetChecked(p.Default >= 0.5)

			check.OnChanged = func(on bool) {
				val := 0.0
				if on {
					val = 1
				}

				s.cb.OnParam(kind, p.Key, val)
			}

			w.check = check
			items = append(items, check)

		default:
			slider := widget.NewSlider(p.Min, p.Max)
			slider.Step = p.Step
			slider.Value = p.Default

			slider.OnChanged = func(v float64) {
				label.SetText(formatParam(p, v))
				s.cb.OnParam(kind, p.Key, v)
			}

			w.slider = slider
			items = append(items, label, slider)
		}

		s.controls[k] = w
	}

	if !def.Live {
		items = append(
			items,
			widget.NewButton("Apply "+def.Label, func() {
				s.cb.OnCommit(kind)
			}),
		)
	}

	// Bottom border.
	items = append(items, widget.NewSeparator())

	// Wrap the section and add a small fixed gap underneath it.
	section := container.NewVBox(items...)

	gap := canvas.NewRectangle(color.Transparent)
	gap.SetMinSize(fyne.NewSize(1, 4))

	return []fyne.CanvasObject{
		section,
		gap,
	}
}

// ResetOp snaps one op's controls back to their defaults without firing
// OnChanged (Fyne's SetValue/SetChecked only update the widgets).
func (s *Sidebar) ResetOp(kind OpKind) {
	for _, p := range registry[kind].Params {
		w, ok := s.controls[ctlKey{kind: kind, key: p.Key}]
		if !ok {
			continue
		}

		switch {
		case w.check != nil:
			w.check.SetChecked(p.Default >= 0.5)

		case w.slider != nil:
			w.slider.SetValue(p.Default)
		}

		w.label.SetText(formatParam(p, p.Default))
	}
}

// Reset snaps every control back to defaults.
func (s *Sidebar) Reset() {
	for _, kind := range registryOrder {
		s.ResetOp(kind)
	}
}

func formatParam(p Param, val float64) string {
	return fmt.Sprintf("%s: %g", p.Label, val)
}

// CanvasObject returns the displayable Fyne object.
func (s *Sidebar) CanvasObject() fyne.CanvasObject {
	return container.NewVScroll(s.container)
}
