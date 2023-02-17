package utils

import (
	"image/color"
)

func ColorLerp(c0, c1 color.Color, t float64) color.Color {
	r0, g0, b0, a0 := c0.RGBA()
	r1, g1, b1, a1 := c1.RGBA()

	return color.RGBA{
		Lerp(r0, r1, t),
		Lerp(g0, g1, t),
		Lerp(b0, b1, t),
		Lerp(a0, a1, t),
	}
}

func Lerp(a, b uint32, t float64) uint8 {
	return uint8(int32(float64(a)*(1.0-t)+float64(b)*t) >> 8)
}
