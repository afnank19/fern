package main

import (
	"fmt"

	"github.com/afnank19/fern/geometric"
	"github.com/afnank19/fern/image"
)

func main() {
	fmt.Println("Hello, Fern! We will be processing images!")

	_, rgbaImg := image.LoadImage("./assets/evelyn.png")
	// _, rgbaImg := image.LoadImage("./assets/samples/sample4x4.png")

	// image.IterateImage(testImage)
	// point.Invert(rgbaImg)
	// point.Brightness(rgbaImg, 40)
	// point.Grayscale(rgbaImg)
	// point.AvgGrayscale(rgbaImg)
	///point.PhotoshopGrayscale(rgbaImg)
	// point.LinearContrast(rgbaImg, 1.5)
	// point.FastGrayscale(rgbaImg)
	// point.SigmoidalContrast(rgbaImg, 0.10)
	// start := time.Now()
	// point.FastSigmoidalContrast(rgbaImg, 8.0)
	// elapsed := time.Since(start)
	// fmt.Println(" Elapsed -", elapsed)
	// point.Threshold(rgbaImg, 50)

	// noise.Uniform(rgbaImg, 20, true)
	// noise.Gaussian(rgbaImg, 20, true)
	// noise.SaltAndPepper(rgbaImg, 0.005)

	// filter.SliceShenanas()
	// image.IterateImage(rgbaImg)
	// image.TestNegAccess(rgbaImg)
	// newImg := filter.BoxBlur(rgbaImg, 1.0)
	// fmt.Println(filter.CalculateKernelSize(rgbaImg, 1.0))
	//
	// start := time.Now()
	// newImg := filter.GaussianBlur(rgbaImg, 0.01)
	// elapsed := time.Since(start)
	// fmt.Println(" Elapsed -", elapsed)

	// newImg := filter.Sharpen(rgbaImg, 0.5)

	// filter.UnsharpMask(rgbaImg, 0.2, 1.5)
	// composite.NaiveBloom(rgbaImg, 0.30, 0.30, 0.8)
	// noise.Gaussian(rgbaImg, 20, true)
	// noise.Uniform(rgbaImg, 5, true)

	// newImg := filter.Downsample2x(rgbaImg)
	// newImg = filter.Downsample2x(newImg)
	// image.SaveBlah()

	geometric.ChromaticAberration(rgbaImg, 20, 0)

	image.SaveImage(rgbaImg, "fern.png", "./assets/saves")

	// start := time.Now()
	// newImg := filter.GaussianBlur(rgbaImg, 1.5)
	// elapsed := time.Since(start)
	// fmt.Println(" NAIVE GAUSSIAN: Elapsed -", elapsed)

	// image.SaveImage(newImg, "naive-gauss.png", "./assets/saves")

	fmt.Println("----------------------------")

	SliceShift()

	// start = time.Now()
	// fastGauss := filter.FastGaussianBlur(rgbaImg, 1.5)
	// elapsed = time.Since(start)
	// fmt.Println(" FAST GAUSSIAN: Elapsed -", elapsed)

	// image.SaveImage(fastGauss, "fast-gauss.png", "./assets/saves")
}

// TODO: remove
func SliceShift() {
	temp := []int{1, 2, 2, 2, 3, 2, 2, 2, 5, 2, 2, 2, 7}
	// slice := []int{1, 2, 2, 2, 3, 2, 2, 2, 5, 2, 2, 2, 7}

	// t := temp[0]
	// temp[0] = 0
	// for i := 0; i < len(temp)-4; i += 4 {
	// 	f := t
	// 	t = temp[i+4]
	// 	temp[i+4] = f
	// }

	pos := 3
	offset := 4
	n := len(temp)

	for i := n-1; i >= offset * pos; i -= offset {
		newVal := temp[i-(offset * pos)]
		temp[i] = newVal
	}

	// if (offset * pos > n) {
	// 	panic("beabadoobee", )
	// }

	for i := 0; i < offset * pos; i += offset {
		temp[i] = 0
	}

	fmt.Println(temp)
	// SliceShiftWithZeroes(slice, 2)
	// SliceShiftEvery4thWithZeroes(slice, 2)
}
