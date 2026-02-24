package point

import (
	"image"
	"runtime"
	"sync"
)

// this works for darkness as well given a negative valuw
func Brightness(img *image.RGBA, value int) {
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// access pixel at (x, y)
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]

			newR := int(r) + (value)
			newG := int(g) + (value)
			newB := int(b) + (value)

			img.Pix[i+0] = clamp(newR)
			img.Pix[i+1] = clamp(newG)
			img.Pix[i+2] = clamp(newB)
		}
	}
}

func BrightnessMultiThread(img *image.RGBA, value int) {
	bounds := img.Bounds()
    height := bounds.Max.Y - bounds.Min.Y
    numWorkers := runtime.NumCPU()
    rowsPerWorker := (height + numWorkers - 1) / numWorkers

    var wg sync.WaitGroup
    for w := 0; w < numWorkers; w++ {
        startY := bounds.Min.Y + w*rowsPerWorker
        endY := min(startY+rowsPerWorker, bounds.Max.Y)
        if startY >= endY {
            break
        }
        wg.Add(1)
        go func(startY, endY int) {
            defer wg.Done()
            for y := startY; y < endY; y++ {
                for x := bounds.Min.X; x < bounds.Max.X; x++ {
                    i := img.PixOffset(x, y)
                    img.Pix[i+0] = clamp(int(img.Pix[i+0]) + value)
                    img.Pix[i+1] = clamp(int(img.Pix[i+1]) + value)
                    img.Pix[i+2] = clamp(int(img.Pix[i+2]) + value)
                }
            }
        }(startY, endY)
    }
    wg.Wait()
}

// may move to a utils later
func clamp(value int) uint8 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return uint8(value)
}
