package utils

import (
	"github.com/fogleman/gg"
)

/* From https://www.arduino.cc/reference/en/language/functions/math/map/ */
func Map(x, originMin, originMax, targetMin, targetMax float64) float64 {
	return (x-originMin)*(targetMax-targetMin)/(originMax-originMin) + targetMin
}

/* Clamp the value between [low, high]*/
func Clamp(f, low, high float64) float64 {
	if f < low {
		return low
	}
	if f > high {
		return high
	}
	return f
}

/* Convert the rotation to coordinates for a linear gradient */
func ToLinerarGradientCoordinates(rotation int) (startX, startY, endX, endY float64) {
	// 0 degrees should be left to right, offset by 45 to do so.
	rotation -= 45

	startPoint := gg.Point{X: 0, Y: 0}
	endPoint := gg.Point{X: 1000, Y: 1000}

	// Now rotate the coordinates to create the desired direction
	rad := gg.Radians(float64(rotation))
	rotationMatrix := RotateAbout(rad, 500, 500, gg.Identity())
	//rotationMatrix.Rotate()
	startPoint.X, startPoint.Y = rotationMatrix.TransformPoint(startPoint.X, startPoint.Y)
	endPoint.X, endPoint.Y = rotationMatrix.TransformPoint(endPoint.X, endPoint.Y)

	return startPoint.X, startPoint.Y, endPoint.X, endPoint.Y
}

// RotateAbout updates the current matrix with a anticlockwise rotation.
// Rotation occurs about the specified point. Angle is specified in radians.
func RotateAbout(angle, x, y float64, matrix gg.Matrix) gg.Matrix {
	return matrix.Translate(x, y).Rotate(angle).Translate(-x, -y)
}
