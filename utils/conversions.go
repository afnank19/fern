package utils

import "math"

func ClampFastFloat(val float64) uint8 {
	// This achieves the same result as the if-statements
	return uint8(math.Max(0, math.Min(255, val)))
}

func FloatToUint8(v float64) uint8 {
    v = math.Max(0, math.Min(1, v)) // clamp to [0,1]
    return uint8(math.Round(v * 255))
}

func SRGBToLinear(c uint8) float64 {
	val := float64(c) / 255.0
	if val <= 0.04045 {
		return val / 12.92
	}
	return math.Pow((val+0.055)/1.055, 2.4)
}
