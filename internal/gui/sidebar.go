package gui

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/afnank19/fern/utils"
)

// Sidebar owns the sliders and notifies the caller via onChange
// whenever the user adjusts a value.
type Sidebar struct {
	container   *fyne.Container
	adjustments Adjustments
	onChange    func(Adjustments)

	brightnessSlider *widget.Slider
	contrastSlider   *widget.Slider
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

	debouncedChange := debounce(10*time.Millisecond, func() {
		s.onChange(s.adjustments)
	})

	s.brightnessSlider.OnChanged = func(v float64) {
		s.adjustments.Brightness = int(v)
		brightnessLabel.SetText(fmt.Sprintf("Brightness: %+d", int(v)))
		debouncedChange()
	}

	contrastLabel := widget.NewLabel("Contrast: 0")
	s.contrastSlider = widget.NewSlider(-10, 10)
	s.contrastSlider.Step = 1
	s.contrastSlider.Value = 0

	s.contrastSlider.OnChanged = func(v float64) {
		s.adjustments.Contrast = int(v)
		contrastLabel.SetText(fmt.Sprintf("Contrast: %+d", int(v)))
		debouncedChange()
	}

	// --- BASIC TAB ---
	basicTab := container.NewVBox(
		widget.NewLabelWithStyle("Basic Adjustments", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		brightnessLabel,
		s.brightnessSlider,
		contrastLabel,
		s.contrastSlider,
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Basic", basicTab),
		// container.NewTabItem("Color", colorTab),
		container.NewTabItem("Effects", s.buildEffectsTab()),
	)

	tabs.SetTabLocation(container.TabLocationTop)

	if s.container == nil {
		s.container = container.NewVBox(tabs)
	} else {
		s.container.Objects = []fyne.CanvasObject{tabs}
		s.container.Refresh()
	}
}

// Reset returns all sliders to their zero position (does not trigger onChange).
func (s *Sidebar) Reset() {
	s.adjustments = Adjustments{}
	s.brightnessSlider.SetValue(0)
	s.build()
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
		timer = time.AfterFunc(d, func() {
			fyne.Do(f)
		})
	}
}

func (s *Sidebar) buildEffectsTab() fyne.CanvasObject {
	intensityLabel := widget.NewLabel("Intensity: 0")

	intensitySlider := widget.NewSlider(0, 1)
	intensitySlider.Step = 0.01
	intensitySlider.Value = 0

	intensitySlider.OnChanged = func(v float64) {
		s.adjustments.BloomAdj.Intensity = v
		intensityLabel.SetText(fmt.Sprintf("Intensity: %+f", v))
		// debouncedChange()
	}

	thresholdLabel := widget.NewLabel("Threshold: 0.95")

	thresholdSlider := widget.NewSlider(0, 1)
	thresholdSlider.Step = 0.01
	thresholdSlider.Value = 0.95
	// Have to initialize the value otherwise it remains at 0
	s.adjustments.BloomAdj.Threshold = thresholdSlider.Value

	thresholdSlider.OnChanged = func(v float64) {
		s.adjustments.BloomAdj.Threshold = v
		thresholdLabel.SetText(fmt.Sprintf("Threshold: %+f", v))
		// debouncedChange()
	}

	blurAmtLabel := widget.NewLabel("Blur Amount: 0")
	blurAmtSlider := widget.NewSlider(0, 1)
	blurAmtSlider.Step = 0.01
	blurAmtSlider.Value = 0

	blurAmtSlider.OnChanged = func(v float64) {
		s.adjustments.BloomAdj.BlurAmt = v
		blurAmtLabel.SetText(fmt.Sprintf("Blur Amount: %+f", v))
	}

	// No idea if this will freeze the UI, may need to update the func a bit to run in another thread
	bloomApplyButton := widget.NewButton("Apply Bloom", func() {
		s.onChange(s.adjustments)
	})

	noiseLabel := widget.NewLabel("Noise: 0")

	noiseSlider := widget.NewSlider(0, 50)
	noiseSlider.Step = 1
	noiseSlider.Value = 0

	debouncedChange := debounce(10*time.Millisecond, func() {
		s.onChange(s.adjustments)
	})

	noiseSlider.OnChanged = func(v float64) {
		s.adjustments.Noise.NoiseAmt = v
		noiseLabel.SetText(fmt.Sprintf("Intensity: %+f", v))
		debouncedChange()
	}

	noisePerChanButton := widget.NewCheck("Per Channel", func(checked bool) {
		s.adjustments.Noise.PerChan = checked
		debouncedChange()
	})

	return container.NewVBox(
		utils.BuildTitleLabel("Bloom"),
		intensityLabel,
		intensitySlider,
		thresholdLabel,
		thresholdSlider,
		blurAmtLabel,
		blurAmtSlider,
		bloomApplyButton,
		utils.BuildTitleLabel("Noise"),
		noiseLabel,
		noiseSlider,
		noisePerChanButton,
	)
}
