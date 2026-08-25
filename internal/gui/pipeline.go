package gui

import (
	"image"

	"github.com/afnank19/fern/composite"
	"github.com/afnank19/fern/noise"
	"github.com/afnank19/fern/point"
)

// OpKind identifies an operation in the registry.
type OpKind string

const (
	OpBrightness OpKind = "brightness"
	OpContrast   OpKind = "contrast"
	OpBloom      OpKind = "bloom"
	OpNoise      OpKind = "noise"
)

// Params holds one operation's parameter values, keyed by name.
type Params map[string]float64

// Op is a single operation with its parameters: an entry of the commit
// stack, or a pending live adjustment.
type Op struct {
	Kind   OpKind
	Params Params
}

// Widget selects the sidebar control used to edit a parameter.
type Widget int

const (
	SliderWidget Widget = iota
	CheckWidget
)

// Param describes one editable parameter of an operation.
type Param struct {
	Key     string
	Label   string
	Min     float64
	Max     float64
	Step    float64
	Default float64
	Widget  Widget
}

// opDef describes an operation: how it is edited and how it is applied.
//
// Live ops are cheap enough to re-render while their sliders move; non-live
// ops are staged and only executed when committed to the stack.
type opDef struct {
	Label    string
	Category string
	Live     bool
	Params   []Param
	Apply    func(img *image.RGBA, p Params)
}

// registryOrder fixes canonical execution order; map iteration is unordered.
var registryOrder = []OpKind{OpBrightness, OpContrast, OpBloom, OpNoise}

var registry = map[OpKind]opDef{
	OpBrightness: {
		Label:    "Brightness",
		Category: "Basic",
		Live:     true,
		Params: []Param{
			{Key: "amount", Label: "Amount", Min: -100, Max: 100, Step: 1},
		},
		Apply: func(img *image.RGBA, p Params) {
			if p["amount"] == 0 {
				return
			}
			point.BrightnessMultiThread(img, int(p["amount"]))
		},
	},
	OpContrast: {
		Label:    "Contrast",
		Category: "Basic",
		Live:     true,
		Params: []Param{
			{Key: "factor", Label: "Factor", Min: -10, Max: 10, Step: 1},
		},
		Apply: func(img *image.RGBA, p Params) {
			point.FastSigmoidalContrast(img, p["factor"])
		},
	},
	OpBloom: {
		Label:    "Bloom",
		Category: "Effects",
		Live:     false,
		Params: []Param{
			{Key: "intensity", Label: "Intensity", Min: 0, Max: 1, Step: 0.01},
			{Key: "threshold", Label: "Threshold", Min: 0, Max: 1, Step: 0.01, Default: 0.95},
			{Key: "blurAmt", Label: "Blur Amount", Min: 0, Max: 1, Step: 0.01},
		},
		Apply: func(img *image.RGBA, p Params) {
			if p["intensity"] <= 0 || p["blurAmt"] <= 0 {
				return
			}
			composite.Bloom(img, p["intensity"], p["threshold"], p["blurAmt"])
		},
	},
	OpNoise: {
		Label:    "Gaussian Noise",
		Category: "Effects",
		Live:     false,
		Params: []Param{
			{Key: "amount", Label: "Amount", Min: 0, Max: 50, Step: 1},
			{Key: "perChan", Label: "Per Channel", Min: 0, Max: 1, Step: 1, Default: 0, Widget: CheckWidget},
		},
		Apply: func(img *image.RGBA, p Params) {
			if p["amount"] <= 0 {
				return
			}
			noise.Gaussian(img, p["amount"], p["perChan"] >= 0.5)
		},
	},
}

// defaultParams returns Params filled with every declared default value.
func defaultParams(kind OpKind) Params {
	p := make(Params)
	for _, param := range registry[kind].Params {
		p[param.Key] = param.Default
	}
	return p
}

// applyOp executes a registered operation, filling any missing parameters
// with their defaults so stored Ops stay forward-compatible.
func applyOp(img *image.RGBA, op Op) {
	def, ok := registry[op.Kind]
	if !ok || def.Apply == nil {
		return
	}

	p := defaultParams(op.Kind)
	for k, v := range op.Params {
		p[k] = v
	}
	def.Apply(img, p)
}

// cloneParams deep-copies a Params map.
func cloneParams(p Params) Params {
	out := make(Params, len(p))
	for k, v := range p {
		out[k] = v
	}
	return out
}

// RenderState is a self-contained snapshot of everything needed to render a
// frame: the committed stack plus live pending values. Snapshots are cloned
// on creation so the renderer never sees mutable shared maps.
type RenderState struct {
	Stack   []Op
	Pending map[OpKind]Params
}

func (s RenderState) clone() RenderState {
	out := RenderState{
		Stack:   make([]Op, len(s.Stack)),
		Pending: make(map[OpKind]Params, len(s.Pending)),
	}
	copy(out.Stack, s.Stack)
	for kind, p := range s.Pending {
		out.Pending[kind] = cloneParams(p)
	}
	return out
}

// execute renders base+state into dst: committed ops first in stack order,
// then live pending ops in canonical registry order.
func execute(dst, base *image.RGBA, st RenderState) {
	copy(dst.Pix, base.Pix)

	for _, op := range st.Stack {
		applyOp(dst, op)
	}
	for _, kind := range registryOrder {
		if registry[kind].Live {
			if p, ok := st.Pending[kind]; ok {
				applyOp(dst, Op{Kind: kind, Params: p})
			}
		}
	}
}
