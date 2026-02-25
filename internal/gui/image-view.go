package gui

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/afnank19/fern/filter"
)

// ImageView holds the displayed canvas and the original unmodified image.
// All processing is applied to a copy, keeping the original pristine.
type ImageView struct {
	canvasImage *canvas.Image
	original    *image.RGBA
	display     *image.RGBA
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

func (v *ImageView) LoadImage(img *image.RGBA) {
	v.original = img
	v.display = v.original

	canvas := fyne.CurrentApp().Driver().CanvasForObject(v.canvasImage)
	if canvas == nil {
		return
	}

	size := v.canvasImage.Size()
	scale := canvas.Scale()

	targetW := int(size.Width * scale)
	targetH := int(size.Height * scale)
	targetPixels := targetW * targetH

	for {
		b := v.display.Bounds()
		w := b.Dx()
		h := b.Dy()

		if w*h <= targetPixels {
			break
		}

		if w <= 1 || h <= 1 {
			break
		}

		v.display = filter.Downsample2x(v.display)
	}

	v.working = image.NewRGBA(v.display.Bounds())
	v.canvasImage.Image = v.working
	v.canvasImage.Refresh()
}

// ApplyAdjustments clones the original and applies the full adjustment pipeline,
// ensuring no quality loss from repeated edits.
func (v *ImageView) ApplyAdjustments(adj Adjustments) {
	if v.original == nil {
		return
	}

	copy(v.working.Pix, v.display.Pix)
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
