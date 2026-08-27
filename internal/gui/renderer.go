package gui

import (
	"sync"
	"sync/atomic"

	"fyne.io/fyne/v2"
)

// renderer renders preview frames on a single background goroutine.
//
// Why one worker rather than one goroutine per request: slider drags fire far
// faster than heavy operations finish. A render per event would run several
// expensive pipelines in parallel (CPU thrash) and complete out of order,
// letting stale frames overwrite newer ones. A single worker plus a
// latest-wins adjustment slot collapses bursts and keeps exactly one render
// in flight.
type renderer struct {
	view *ImageView

	kick chan struct{} // size 1: non-blocking send signals "a render is wanted"

	mu sync.Mutex // guards st
	st RenderState

	gen atomic.Uint64 // bumped when the base image changes; invalidates in-flight frames
}

func newRenderer(view *ImageView) *renderer {
	r := &renderer{
		view: view,
		kick: make(chan struct{}, 1),
	}
	go r.loop()
	return r
}

// request schedules a re-render with the given state. Calls made while a
// render is running collapse into one follow-up render using the latest
// state.
func (r *renderer) request(st RenderState) {
	r.mu.Lock()
	r.st = st
	r.mu.Unlock()

	select {
	case r.kick <- struct{}{}:
	default:
		// Already queued or rendering; the worker re-reads st when done.
	}
}

// invalidate drops any frame still being rendered from the old base image.
func (r *renderer) invalidate() {
	r.gen.Add(1)
}

func (r *renderer) loop() {
	for range r.kick {
		r.mu.Lock()
		st := r.st
		r.mu.Unlock()

		r.renderFrame(st)
	}
}

func (r *renderer) renderFrame(st RenderState) {
	v := r.view

	v.mu.Lock()
	base := v.display
	scratch := v.scratch
	v.mu.Unlock()
	gen := r.gen.Load()

	if base == nil || scratch == nil || len(base.Pix) != len(scratch.Pix) {
		return
	}

	fyne.Do(func() { v.showBusy() })

	execute(scratch, base, st)

	fyne.Do(func() {
		if r.gen.Load() != gen {
			return // a new image was loaded mid-render; drop the stale frame
		}

		v.mu.Lock()
		v.working, v.scratch = v.scratch, v.working
		img := v.working
		v.mu.Unlock()

		v.canvasImage.Image = img
		v.canvasImage.Refresh()
		v.hideBusy()
	})
}
