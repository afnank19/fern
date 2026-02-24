package gui

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// ImageView holds the displayed canvas and the original unmodified image.
// All processing is applied to a copy, keeping the original pristine.
type ImageView struct {
	canvasImage *canvas.Image
	original    *image.RGBA
	working     *image.RGBA
}

func NewImageView() *ImageView {
	img := canvas.NewImageFromImage(nil)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(400, 400))

	return &ImageView{
		canvasImage: img,
	}
}

// LoadImage stores the original and displays it unmodified.
func (v *ImageView) LoadImage(img *image.RGBA) {
	v.original = img
	v.working = image.NewRGBA(img.Bounds()) // allocate once
	v.canvasImage.Image = v.working
	v.canvasImage.Refresh()
}

// ApplyAdjustments clones the original and applies the full adjustment pipeline,
// ensuring no quality loss from repeated edits.
func (v *ImageView) ApplyAdjustments(adj Adjustments) {
	if v.original == nil {
		return
	}

	copy(v.working.Pix, v.original.Pix)
	applyPipeline(v.working, adj)

	v.canvasImage.Image = v.working
	v.canvasImage.Refresh()
}

// CanvasObject returns the displayable Fyne object.
func (v *ImageView) CanvasObject() fyne.CanvasObject {
	return v.canvasImage
}

// cloneRGBA returns a deep copy of an image.RGBA.
func cloneRGBA(src *image.RGBA) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	copy(dst.Pix, src.Pix)
	return dst
}
