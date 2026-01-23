package geometric

import (
	"image"
)

var offset = 4

// var strength = 4

// Pass in 0 for Red/Cyan, 1 for Blue/Yellow, 2 for Green Magenta
func ChromaticAberration(img *image.RGBA, strength int, fringeType int) {
	// bounds := img.Bounds()
	// y := bounds.Min.Y
	// x := bounds.Min.X

	// i := img.PixOffset(x, y)

	n := len(img.Pix)

	for idx := n - offset + fringeType; idx >= offset*strength; idx -= offset {
		newVal := img.Pix[idx-(offset*strength)]
		img.Pix[idx] = newVal
	}

	for idx := 0; idx < offset*strength; idx += offset {
		img.Pix[idx] = 0
	}
}
