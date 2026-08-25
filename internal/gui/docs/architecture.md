# Architecture: data flow and editing model

## Big picture

```
 Sidebar (auto-generated controls)
     │  OnParam / OnCommit / OnUndo   ← plain function calls, no image logic
     ▼
 ImageView ──── stack []Op          (committed effects, ordered)
     │           pending map        (slider values not yet committed)
     │
     │  requestRender() → snapshot (deep copy of stack+pending)
     ▼
 renderer (background goroutine) ── one frame at a time, latest wins
     │  execute(scratch, display, state)   ← the only place pixels change
     ▼
 fyne.Do(...) on the UI thread: swap scratch↔working, canvas.Refresh()
```

Two rules explain most of the design:

1. **The original image is sacred.** Nothing ever mutates it. Previews run on
   a downscaled copy; export runs on a fresh full-res copy.
2. **UI code never touches pixels.** It only changes *state* (`stack`,
   `pending`) and asks for a re-render.

## The five buffers (`image-view.go`)

| Field | Resolution | Purpose |
|---|---|---|
| `original` | full | The loaded photo. Never modified, ever. |
| `display` | downscaled | The preview base. Repeatedly halved until it fits the canvas, so effects stay fast. |
| `working` | downscaled | Currently shown on screen. |
| `scratch` | downscaled | Private drawing surface for the render worker; never attached to the canvas (see [concurrency.md](concurrency.md)). |
| `render` | — | Not used for storage anymore: `PrepareExport()` builds a fresh copy per export. |

**Why separate buffers at all?** An earlier version assigned
`display = original` ("just another name for the same memory"). Export then
wrote its result through that shared reference into the *original*, so every
export permanently baked the edits into the source photo. Deep copies
(`CopyRGBA`) make each buffer independent.

**Why downscaled previews?** Bloom can take hundreds of milliseconds at full
resolution but only a few at preview size. Halving repeatedly
(`filter.Downsample2x`) is also cheaper than resizing once to exact size.

## Editing model: stack vs pending

- **Live ops** (`Live: true`, e.g. Brightness): cheap point operations.
  Moving their slider updates `pending` and triggers an async preview
  immediately. They are never committed — they are recomputed from scratch on
  every frame, so there is no quality loss and no stacking bug.
- **Staged ops** (`Live: false`, e.g. Bloom): expensive effects. Sliders only
  update `pending`; nothing renders until you press **Apply**, which moves
  the values onto `stack`. Undo pops the stack. This was chosen over
  "live-preview everything" because dragging a slider would otherwise rerun a
  multi-pass blur dozens of times.

A preview frame is always:

```
copy of display  →  apply stack in order  →  apply live pending values
```

Execution order is deterministic: stack in commit order, then live pending
ops in `registryOrder` sequence. That ordering lives in exactly one place
(`execute` in `pipeline.go`).

## Why a registry?

Before, every slider was hand-coded twice: once as a widget in the sidebar,
once as an `if` in the pipeline. They could drift apart silently. Now
`registry` (`pipeline.go`) is the single source of truth describing each op's
name, tab category, parameter sliders, defaults, and apply function. The
sidebar walks the registry to build itself, and `applyOp` walks it to execute.

Consequences:

- Adding an effect = adding registry data (see [adding-an-op.md](adding-an-op.md)).
- Widget code and pixel code cannot disagree about parameters.
- Defaults live next to ranges, and missing values are auto-filled
  (`defaultParams`), so stored ops keep working when new params are added.

## Export path

`PrepareExport()` snapshots the current state, deep-copies `original` at full
resolution, runs the same `execute()` on it, and hands the result to the save
dialog. Because the input is always a pristine copy, exporting twice gives
the same file — exports are idempotent by construction.

(Export currently runs synchronously on the UI thread. Acceptable today;
making it async with a progress dialog is a known follow-up.)
