package gui

import (
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
//
// Editing model: committed ops live in stack (heavy effects pushed via
// Commit), cheap live tweaks accumulate in pending and are re-rendered on
// every change. A preview frame is always base + stack + live pending.
type ImageView struct {
	canvasImage *canvas.Image

	mu       sync.Mutex  // guards the pointer/state fields below, not pixel data
	original *image.RGBA // pristine full-res source, never mutated
	display  *image.RGBA // downscaled base for previews
	working  *image.RGBA // buffer currently attached to the canvas
	scratch  *image.RGBA // offscreen buffer only the render worker writes to
	render   *image.RGBA // full-res export target

	stack   []Op              // committed operations, execution order
	pending map[OpKind]Params // uncommitted parameter values by op kind

	rend *renderer
}

func NewImageView() *ImageView {
	img := canvas.NewImageFromImage(nil)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(400, 400))

	v := &ImageView{
		canvasImage: img,
		pending:     make(map[OpKind]Params),
	}
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
	// A new image starts with a clean editing state.
	v.stack = nil
	v.pending = make(map[OpKind]Params)
	v.mu.Unlock()

	// Frames still rendering from the old image are now stale.
	v.rend.invalidate()

	v.canvasImage.Image = v.working
	v.canvasImage.Refresh()
}

// SetParam updates a pending parameter value. Live ops trigger an async
// preview re-render; staged ops only update state until committed.
func (v *ImageView) SetParam(kind OpKind, key string, val float64) {
	def, ok := registry[kind]
	if !ok || v.original == nil {
		return
	}

	v.mu.Lock()
	if v.pending[kind] == nil {
		v.pending[kind] = defaultParams(kind)
	}
	v.pending[kind][key] = val
	v.mu.Unlock()

	if def.Live {
		v.requestRender()
	}
}

// Commit moves an op's pending parameters onto the stack as a permanent
// step and re-renders. Used by Apply buttons for heavy effects.
func (v *ImageView) Commit(kind OpKind) {
	v.mu.Lock()
	p, ok := v.pending[kind]
	delete(v.pending, kind)
	if ok && len(p) > 0 {
		v.stack = append(v.stack, Op{Kind: kind, Params: cloneParams(p)})
	}
	v.mu.Unlock()

	v.requestRender()
}

// Undo removes the most recently committed operation and re-renders.
func (v *ImageView) Undo() {
	v.mu.Lock()
	if n := len(v.stack); n > 0 {
		v.stack = v.stack[:n-1]
	}
	v.mu.Unlock()

	v.requestRender()
}

// Reset clears all editing state (used when a new image is loaded).
func (v *ImageView) Reset() {
	v.mu.Lock()
	v.stack = nil
	v.pending = make(map[OpKind]Params)
	v.mu.Unlock()
}

// snapshot returns a fully cloned view of the current editing state, safe to
// hand to the render worker or hold beyond the lock.
func (v *ImageView) snapshot() RenderState {
	v.mu.Lock()
	defer v.mu.Unlock()
	return RenderState{
		Stack:   append([]Op(nil), v.stack...),
		Pending: v.pending,
	}.clone()
}

// requestRender schedules an async preview of base+stack+pending on the
// downscaled buffers, so slider drags never block the UI.
func (v *ImageView) requestRender() {
	if v.original == nil {
		return
	}
	v.rend.request(v.snapshot())
}

// PrepareExport executes the full state at original resolution into a fresh
// buffer and returns it. Synchronous by design for now; Phase 4 moves this
// off the UI thread with a progress dialog.
func (v *ImageView) PrepareExport() *image.RGBA {
	st := v.snapshot()

	v.mu.Lock()
	base := v.original
	v.mu.Unlock()

	if base == nil {
		return nil
	}
	out := fernimage.CopyRGBA(base)
	execute(out, base, st)
	return out
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
