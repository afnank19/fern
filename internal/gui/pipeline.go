package gui

import (
	"image"

	"github.com/afnank19/fern/composite"
	"github.com/afnank19/fern/noise"
	"github.com/afnank19/fern/point"
)

// Adjustments holds the current values of all user-controlled parameters.
// Adding a new control is as simple as adding a field here and a step in applyPipeline.
type Adjustments struct {
	Brightness int
	Contrast   int
	// Saturation int
	BloomAdj BloomAdjustments
	Noise    NoiseAdjustments
}

type BloomAdjustments struct {
	Intensity float64
	Threshold float64
	BlurAmt   float64
}

type NoiseAdjustments struct {
	NoiseAmt float64
	PerChan  bool
}

// applyPipeline runs every adjustment in a deterministic order.
// Operations are intentionally sequenced (e.g. brightness before contrast).
func applyPipeline(img *image.RGBA, adj Adjustments) {
	if adj.Brightness != 0 {
		point.BrightnessMultiThread(img, adj.Brightness)
	}
	if adj.Contrast != 0 {
		point.FastSigmoidalContrast(img, float64(adj.Contrast))
	}
	if adj.BloomAdj.Intensity > 0 && adj.BloomAdj.BlurAmt > 0 {
		composite.Bloom(img, adj.BloomAdj.Intensity, adj.BloomAdj.Threshold, adj.BloomAdj.BlurAmt)
	}
	if adj.Noise.NoiseAmt > 0 {
		noise.Gaussian(img, adj.Noise.NoiseAmt, adj.Noise.PerChan)
	}
}
