# Fern GUI — Start Here

This package is the Fyne-based image editor UI. Almost everything it does is
driven by one idea:

> **The UI is generated from a registry of operations, and every edit is just
> data — never direct pixel manipulation from the UI layer.**

Because of that, adding a new effect to the app usually means adding an entry
to a table in `pipeline.go` and nothing else.

## The files (read in this order)

| File | What it does | Doc |
|---|---|---|
| `gui.go` | Assembles the window: top bar, image view, sidebar. ~70 lines of wiring. | |
| `pipeline.go` | **The heart**: op types, the registry (every effect + its sliders), and `execute()` which renders a state into pixels. | [architecture.md](docs/architecture.md) |
| `image-view.go` | Owns all image buffers and editing state (op stack, pending values). Exposes `SetParam / Commit / Undo / Reset`. | [architecture.md](docs/architecture.md) |
| `renderer.go` | One background goroutine that renders preview frames without freezing the UI. | [concurrency.md](docs/concurrency.md) |
| `sidebar.go` | Builds sliders/checkboxes/buttons *automatically* from the registry. Holds no image state. | [architecture.md](docs/architecture.md) |
| `theme.go` | Custom dark theme: palette constants, forced monospace font, square corners. | |

## How to add a new effect

See **[adding-an-op.md](docs/adding-an-op.md)** — it's 3 small edits in
`pipeline.go`, using Chromatic Aberration as a real worked example that is
already in the codebase.

## Mini-glossary (Go terms used in these docs)

- **Pointer (`*T`)**: a reference to a value. Two pointers can reference the
  *same* memory — mutating through one is visible through the other. This is
  why we deep-copy buffers instead of assigning them.
- **Buffer**: here, an `image.RGBA` — an in-memory bitmap (one byte per R, G,
  B, A channel per pixel) plus helper methods.
- **Goroutine**: Go's unit of concurrent execution, like a lightweight thread.
  `go f()` runs `f` "at the same time" as the caller.
- **Mutex (`sync.Mutex`)**: a lock. Only one goroutine can hold it at a time;
  used so two goroutines never touch the same field simultaneously.
- **Channel**: a pipe between goroutines. Sending/receiving also guarantees
  the receiver sees everything the sender did before sending.
- **Deep copy / snapshot**: copying a structure *and* everything it points
  to, so the copy can change without affecting the original.
- **Registry**: a lookup table (`map[OpKind]opDef`) describing each effect:
  its name, its sliders, and how to apply it.
