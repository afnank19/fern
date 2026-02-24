package gui

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Sidebar owns the sliders and notifies the caller via onChange
// whenever the user adjusts a value.
type Sidebar struct {
	container    *fyne.Container
	adjustments  Adjustments
	onChange     func(Adjustments)

	brightnessSlider *widget.Slider
}

func NewSidebar(onChange func(Adjustments)) *Sidebar {
	s := &Sidebar{onChange: onChange}
	s.build()
	return s
}

func (s *Sidebar) build() {
	brightnessLabel := widget.NewLabel("Brightness: 0")

	s.brightnessSlider = widget.NewSlider(-100, 100)
	s.brightnessSlider.Step = 1
	s.brightnessSlider.Value = 0

	debouncedChange := debounce(40*time.Millisecond, func() {
    	s.onChange(s.adjustments)
	})

	s.brightnessSlider.OnChanged = func(v float64) {
		s.adjustments.Brightness = int(v)
		brightnessLabel.SetText(fmt.Sprintf("Brightness: %+d", int(v)))
		// s.onChange(s.adjustments)
		debouncedChange()
	}

	s.container = container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Adjustments", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		brightnessLabel,
		s.brightnessSlider,
		// add more controls here as you implement them
	)
}

// Reset returns all sliders to their zero position (does not trigger onChange).
func (s *Sidebar) Reset() {
	s.adjustments = Adjustments{}
	s.brightnessSlider.SetValue(0)
}

// CanvasObject returns the displayable Fyne object.
func (s *Sidebar) CanvasObject() fyne.CanvasObject {
	return container.NewVScroll(s.container)
}

func debounce(d time.Duration, f func()) func() {
    var mu sync.Mutex
    var timer *time.Timer
    return func() {
        mu.Lock()
        defer mu.Unlock()
        if timer != nil {
            timer.Stop()
        }
        timer = time.AfterFunc(d, f)
    }
}
