# Adding a new operation

The whole recipe is edits to **one file**: `pipeline.go`. The sidebar,
preview renderer, commit/undo, and export pick it up automatically — no
changes needed in `image-view.go`, `sidebar.go`, `renderer.go`, or `gui.go`.

## Step 0: check your function's shape

Registry apply functions receive an image they may **mutate in place**:

```go
Apply: func(img *image.RGBA, p Params) { ... }
```

Most of this repo's ops (point/, noise/, geometric/, composite/) already work
this way. Functions that *return a new image* instead (several in filter/)
need a tiny adapter:

```go
Apply: func(img *image.RGBA, p Params) {
    out := filter.BoxBlur(img, p["amount"])
    copy(img.Pix, out.Pix) // bring the result back into the working buffer
},
```

## Steps 1–3

**1. Declare the kind** with the other constants:

```go
OpChromaticAberration OpKind = "chromatic-aberration"
```

**2. Add the registry entry.** Real example already in the codebase
(`pipeline.go`):

```go
OpChromaticAberration: {
    Label:    "Chromatic Aberration",
    Category: "Effects",        // tab name; first use creates the tab
    Live:     false,            // staged until Apply (see below)
    Params: []Param{
        {Key: "strength", Label: "Strength", Min: 1, Max: 10, Step: 1},
        {Key: "fringeType", Label: "Fringe Type", Min: 0, Max: 2, Step: 1},
    },
    Apply: func(img *image.RGBA, p Params) {
        if p["strength"] < 1 {
            return // guard: unset/invalid params must be a safe no-op
        }
        geometric.ChromaticAberration(img, int(p["strength"]), int(p["fringeType"]))
    },
},
```

**3. Slot it into `registryOrder`.** This list fixes execution order *and*
sidebar ordering:

```go
[]OpKind{OpBrightness, OpContrast, OpBloom, OpChromaticAberration, OpNoise}
```

## Live or staged?

| | `Live: true` | `Live: false` |
|---|---|---|
| UI | slider only | sliders + auto-generated "Apply" button |
| Preview | re-renders while dragging | renders only when committed |
| Cost budget | single fast pass over pixels | anything slower (multi-pass blur etc.) |

When unsure, start `false` and flip later; nothing else changes.

**Always include the zero-guard in Apply.** Missing params are auto-filled
with defaults (`defaultParams`), and numeric defaults are 0 — without a guard,
a half-configured op could run with meaningless values instead of doing
nothing.

## Parameter widgets

- Default is a slider: set `Min/Max/Step/Default`.
- Checkbox: add `Widget: CheckWidget` (value 0/1). Convert back with
  `p["key"] >= 0.5`.
- Dropdowns don't exist yet. To add one (e.g. Grayscale variants):
  extend the `Widget` enum with `SelectWidget`, give `Param` an
  `Options []string` field, render a `widget.Select` in `sidebar.go`
  `buildOpControls`, map the chosen index back via `cb.OnParam(..., float64(idx))`,
  and restore defaults in `ResetOp`.

## Verify

```sh
go build ./... && go vet ./internal/gui/ && gofmt -l internal/gui/
```

Then manually: load an image → move the new controls → Apply (if staged) →
Undo → confirm preview updates → Export and reopen the file to confirm the
effect is baked in at full resolution.
