# Concurrency: how the renderer stays responsive and race-free

One rule governs everything:

> **Fyne runs the UI on a single goroutine. Widgets may only be touched from
> that goroutine.** Heavy work goes to a background goroutine; results come
> back with `fyne.Do(func() { ... })`, which queues code onto the UI
> goroutine.

If you run slow pixel work inside a button callback, the whole app freezes —
the UI goroutine is busy and can't process input or repaint. That's why the
renderer exists.

## Why not just spawn one goroutine per slider event?

Three failures, each alone disqualifying:

1. Widget updates off the UI goroutine race Fyne's painter → crashes.
2. A drag fires dozens of events → dozens of *parallel* Bloom renders → CPU
   thrash, everything slower than being synchronous.
3. Goroutines finish out of order: a frame for "slider=3" can land after
   "slider=9", making the preview flicker backwards.

The fix is not smarter scheduling of many renders — it is having **at most
one render in flight, always rendering the newest state**.

## The four mechanisms (`renderer.go`)

### 1. One worker goroutine
`newRenderer` starts `loop()` once; it sleeps until signaled, renders one
frame, repeats. Parallelism is capped at exactly one render.

### 2. Latest-wins slot + size-1 channel (coalescing)
```go
request(st):   store st in r.st; try to drop a token into kick (never blocks)
loop():        wait for token; read r.st (the LATEST); render it
```
We never queue values — only a boolean-like "something changed" signal. Ten
slider events during one render collapse into a single follow-up render of
the final position. This replaced an old timer-based debounce entirely.

### 3. Generation counter (invalidation)
Scenario the worker model can't solve alone: a render of photo **A** is
mid-flight when the user opens photo **B**. Buffers are replaced; publishing
the finished A-frame would paint the old photo over the new one.

Fix: `gen` is an atomic counter bumped by `LoadImage`. The worker snapshots
it before rendering and re-checks inside its publish callback; on mismatch
the frame is silently dropped. One atomic add + load — the cheapest possible
"is this still relevant?" check.

### 4. Double buffering (why pixels need no locks)
Fyne reads the buffer attached to `canvasImage` while painting. If the worker
wrote those pixels directly, that's a data race (torn frames, undefined
behavior). Instead:

- Worker paints only into `scratch` — a buffer the canvas has never seen.
- Publishing is a **pointer swap** done inside `fyne.Do`, i.e. on the UI
  goroutine itself. Since painting also happens there, a swap cannot
  interleave with a paint: every repaint sees either the fully-old or
  fully-new buffer.
- Therefore pixel data needs no locking at all; the mutex guards only the
  few pointer fields shared between `LoadImage`, the worker, and the swap.

## Locking rules used throughout this package

- Guard *small metadata* (pointers, maps of settings), never big pixel arrays.
- Hold a lock just long enough to copy what you need out; **never across
  heavy computation**.
- Anything handed to another goroutine is deep-copied first
  (`snapshot()` → `RenderState.clone()`), so ownership is unambiguous.

## Checklist for future async features

1. Identify work that takes >~10ms; move it off the UI goroutine.
2. Decide the result policy: previews want latest-wins (reuse `renderer`);
   exports/logs want every result.
3. Publish UI changes via `fyne.Do`.
4. Add an invalidation epoch if some change makes in-flight results obsolete.
5. Run `go build -race` and abuse the sliders — the race detector finds what
   eyes miss.
