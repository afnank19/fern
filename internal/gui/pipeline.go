package gui

import (
	"image"

	"github.com/afnank19/fern/point"
)

// Adjustments holds the current values of all user-controlled parameters.
// Adding a new control is as simple as adding a field here and a step in applyPipeline.
type Adjustments struct {
	Brightness int
	// Contrast   int
	// Saturation int
}

// applyPipeline runs every adjustment in a deterministic order.
// Operations are intentionally sequenced (e.g. brightness before contrast).
func applyPipeline(img *image.RGBA, adj Adjustments) {
	if adj.Brightness != 0 {
		point.BrightnessMultiThread(img, adj.Brightness)
	}
	// processing.Contrast(img, adj.Contrast)
}
