package gui

import (
	"fmt"
	"image"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/afnank19/fern/filter"
	fernimage "github.com/afnank19/fern/image"
)

// ImageView holds the displayed canvas and the buffers the GUI works on.
// original/render stay at full resolution; display/working/scratch are
// downscaled preview buffers.
type ImageView struct {
	canvasImage *canvas.Image

	mu       sync.Mutex  // guards the pointer fields below, not pixel data
	original *image.RGBA // pristine full-res source, never mutated
	display  *image.RGBA // downscaled base for previews
	working  *image.RGBA // buffer currently attached to the canvas
	scratch  *image.RGBA // offscreen buffer only the render worker writes to
	render   *image.RGBA // full-res export target

	rend *renderer
}

func NewImageView() *ImageView {
	img := canvas.NewImageFromImage(nil)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(400, 400))

	v := &ImageView{canvasImage: img}
	v.rend = newRenderer(v)
	return v
}

func (v *ImageView) LoadImage(img *image.RGBA) {
	display := fernimage.CopyRGBA(img)
	render := fernimage.CopyRGBA(img)

	// Downscale until the preview fits the canvas; processing runs against
	// this smaller buffer so previews stay responsive.
	if canvas := fyne.CurrentApp().Driver().CanvasForObject(v.canvasImage); canvas != nil {
		size := v.canvasImage.Size()
		targetPixels := int(size.Width*canvas.Scale()) * int(size.Height*canvas.Scale())

		for {
			b := display.Bounds()
			if b.Dx()*b.Dy() <= targetPixels || b.Dx() <= 1 || b.Dy() <= 1 {
				break
			}
			display = filter.Downsample2x(display)
		}
	}

	v.mu.Lock()
	v.original = img
	v.render = render
	v.display = display
	// Seed working so the new image shows immediately instead of an empty
	// (black) buffer, and give the worker a matching scratch buffer.
	v.working = fernimage.CopyRGBA(display)
	v.scratch = fernimage.CopyRGBA(display)
	v.mu.Unlock()

	// Frames still rendering from the old image are now stale.
	v.rend.invalidate()

	v.canvasImage.Image = v.working
	v.canvasImage.Refresh()
}

// ApplyAdjustments schedules an async preview re-render of the full
// adjustment pipeline over the downscaled base, so slider drags never block
// the UI. Results land on the canvas via the renderer's fyne.Do callback.
func (v *ImageView) ApplyAdjustments(adj Adjustments) {
	if v.original == nil {
		return
	}

	v.rend.request(adj)
}

func (v *ImageView) ApplyAdjustmentsOnImage(adj Adjustments) {
	fmt.Println("Applying Adj on original image", adj)
	// render owns its own pixels, so this is a real reset from the pristine
	// original before applying the pipeline at full resolution. Repeated
	// exports therefore never stack adjustments.
	copy(v.render.Pix, v.original.Pix)
	applyPipeline(v.render, adj)
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
