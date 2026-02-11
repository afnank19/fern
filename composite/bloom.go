package composite

import (
	"fmt"
	"image"
	"math"

	"github.com/afnank19/fern/filter"
	localImg "github.com/afnank19/fern/image"
	"github.com/afnank19/fern/point"
	"github.com/afnank19/fern/utils"
)

const samplePasses = 8

func NaiveBloom(img *image.RGBA, intensity ,threshold, blurAmt float64) {
	bounds := img.Bounds()

	brightPass := image.NewRGBA(bounds)

	// Create an image that has the bright parts based on the threshold
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// access pixel at (x, y)
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]
			a := img.Pix[i+3]

			lum := point.LuminancePhotoshop(r, g, b)
			// IMPORTANT: convert either on to the other properly, not just type conv
			if (lum < utils.FloatToUint8(threshold)) {
				// fmt.Println("lum", lum, "thre", utils.FloatToUint8(threshold))
				// Softer bloom threshold
				t := float64(lum)/255.0
				w := math.Max(0, (t-threshold)/(1-threshold))

				r = uint8(float64(r) * w)
				g = uint8(float64(g) * w)
				b = uint8(float64(b) * w)
			}

			// brightIndex := brightPass.PixOffset(x, y)
			brightPass.Pix[i + 0] = r
			brightPass.Pix[i + 1] = g
			brightPass.Pix[i + 2] = b
			brightPass.Pix[i + 3] = a
		}
	}


	localImg.SaveImage(brightPass, "brightpass.png", "./assets/saves")
	blurredBright := filter.FastGaussianBlur(brightPass, blurAmt)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// access pixel at (x, y)
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]

			img.Pix[i+0] = addBloomToPixel(r, intensity, blurredBright.Pix[i+0])
			img.Pix[i+1] = addBloomToPixel(g, intensity, blurredBright.Pix[i+1])
			img.Pix[i+2] = addBloomToPixel(b, intensity, blurredBright.Pix[i+2])
		}
	}

}

func addBloomToPixel(original uint8, intensity float64, blurredHighlights uint8) uint8 {
	return utils.ClampFastFloat(float64(original) + float64(blurredHighlights) * intensity)
}


func Bloom(img *image.RGBA, intensity ,threshold, blurAmt float64) {
	bounds := img.Bounds()

	brightPass := image.NewRGBA(bounds)

	// Create an image that has the bright parts based on the threshold
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// access pixel at (x, y)
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]
			a := img.Pix[i+3]

			lum := point.LuminancePhotoshop(r, g, b)
			// IMPORTANT: convert either on to the other properly, not just type conv
			if (lum < utils.FloatToUint8(threshold)) {
				// fmt.Println("lum", lum, "thre", utils.FloatToUint8(threshold))
				// Softer bloom threshold
				t := float64(lum)/255.0
				w := math.Max(0, (t-threshold)/(1-threshold))

				r = uint8(float64(r) * w)
				g = uint8(float64(g) * w)
				b = uint8(float64(b) * w)
			}

			// brightIndex := brightPass.PixOffset(x, y)
			brightPass.Pix[i + 0] = r
			brightPass.Pix[i + 1] = g
			brightPass.Pix[i + 2] = b
			brightPass.Pix[i + 3] = a
		}
	}

	localImg.SaveImage(brightPass, "bloom-brightpass.png", "./assets/bloom")

	var downsampledBlurredPasses []*image.RGBA
	var levelWeights []float64
	levelWeights = append(levelWeights, 0.5)
	blurKernel := filter.GetKernelForBloom(img, blurAmt)

	activeBrightPass := localImg.CopyRGBA(brightPass)
	for i := 0; i < samplePasses; i++ {
		downsampledImg := filter.Downsample2x(activeBrightPass)

		// downsampledBlurred := filter.FastGaussianBlur(downsampledImg, blurAmt)
		downsampledBlurred := filter.FastGaussianBlurWithKernel(downsampledImg, blurKernel)
		localImg.SaveImage(downsampledBlurred, fmt.Sprintf("down-sample-%dx.png",i*2), "./assets/bloom")
		downsampledBlurredPasses = append(downsampledBlurredPasses, downsampledBlurred)

		activeBrightPass = localImg.CopyRGBA(downsampledBlurred)
		levelWeights = append(levelWeights, levelWeights[i] / 2.0)
	}

	fmt.Println(len(levelWeights), levelWeights)
	fmt.Println(len(downsampledBlurredPasses))
	for idx := len(downsampledBlurredPasses)-1; idx > 0; idx-- {
		downsampledImg := localImg.CopyRGBA(downsampledBlurredPasses[idx])
		upsampledImg := filter.Upsample2x(downsampledImg)
		localImg.SaveImage(upsampledImg, fmt.Sprintf("up-sample-pre-%dx.png",idx*2), "./assets/bloom")
		nextDSI := downsampledBlurredPasses[idx-1]

		minBounds := getMinimumBounds(nextDSI, upsampledImg)
		for y := minBounds.Min.Y; y < minBounds.Max.Y; y++ {
			for x := minBounds.Min.X; x < minBounds.Max.X; x++ {
				// access pixel at (x, y)
				i := nextDSI.PixOffset(x, y)
				r := nextDSI.Pix[i+0]
				g := nextDSI.Pix[i+1]
				b := nextDSI.Pix[i+2]

				iUpsampled := upsampledImg.PixOffset(x, y)
				rUp := upsampledImg.Pix[iUpsampled+0]
				gUp := upsampledImg.Pix[iUpsampled+1]
				bUp := upsampledImg.Pix[iUpsampled+2]

				nextDSI.Pix[i+0] = saturatingWeightedSum(r, rUp, levelWeights[idx])
				nextDSI.Pix[i+1] = saturatingWeightedSum(g, gUp, levelWeights[idx])
				nextDSI.Pix[i+2] = saturatingWeightedSum(b, bUp, levelWeights[idx])
				nextDSI.Pix[i+3] = 255
			}
		}

		localImg.SaveImage(nextDSI, fmt.Sprintf("up-sample-%dx.png",idx*2), "./assets/bloom")
	}

	// The first image in this array is 1 Downsample2x lower than our original image
	// Original -> 6x6 so this will be 3x3, we need to upsample it again
	temp := localImg.CopyRGBA(downsampledBlurredPasses[0])
	localImg.SaveImage(temp ,"temp-downsample-0.png", "./assets/bloom")
	finalUpsample := filter.Upsample2x(temp)
	localImg.SaveImage(finalUpsample,"up-sample-0x.png", "./assets/bloom")

	finalMinBounds := getMinimumBounds(img, finalUpsample)
	for y := finalMinBounds.Min.Y; y < finalMinBounds.Max.Y; y++ {
		for x := finalMinBounds.Min.X; x < finalMinBounds.Max.X; x++ {
			// access pixel at (x, y)
			i := img.PixOffset(x, y)
			r := img.Pix[i+0]
			g := img.Pix[i+1]
			b := img.Pix[i+2]

			iUpsampled := finalUpsample.PixOffset(x, y)
			rUp := finalUpsample.Pix[iUpsampled+0]
			gUp := finalUpsample.Pix[iUpsampled+1]
			bUp := finalUpsample.Pix[iUpsampled+2]

			img.Pix[i+0] = addBloomToPixel(r, intensity, rUp)
			img.Pix[i+1] = addBloomToPixel(g, intensity, gUp)
			img.Pix[i+2] = addBloomToPixel(b, intensity, bUp)
			img.Pix[i+3] = 255
		}
	}
}

func saturatingSum(a, b uint8) uint8 {
	sum := int(a) + int(b)
	if sum > 255 {
		return 255
	}
	return uint8(sum)
}

func saturatingWeightedSum(a, b uint8, weight float64) uint8 {
    // Multiply b by weight, then add to a
    // We use float64 for the intermediate math to maintain precision
    sum := float64(a) + (float64(b) * weight)

    if sum > 255 {
        return 255
    }
    if sum < 0 {
        return 0 // Handle negative weights just in case
    }

    return uint8(sum)
}

func getMinimumBounds(img1, img2 *image.RGBA) image.Rectangle {
    // Intersect returns the largest rectangle contained by both rectangles
    fmt.Println(img1.Bounds())
    fmt.Println(img2.Bounds())
    fmt.Println()
    return img1.Bounds().Intersect(img2.Bounds())
}
