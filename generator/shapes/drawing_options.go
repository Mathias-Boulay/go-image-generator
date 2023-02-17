package shapes

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/mathias-boulay/generator/utils"
)

/* Struct representing parameters for the DrawStepped Blob method */
type SteppedDrawingOptions struct {
	/* Internal stuff: used for the colorAt implementation. Leaky abstraction, sadly */
	currentStep int

	/* How many steps will be executed */
	Steps int

	/* The difference in scale between steps */
	ScaleStep float64

	/* The difference in rotation between steps */
	RotationStep int

	/* The difference in elevation between steps */
	ElevationStep uint32

	/* The difference in coordinates between steps
	* Note that the coordinates interpretation will differ according to Blob settings
	 */
	TranslateStep gg.Point

	/* Start and end patterns, lerping from one to another */
	StartPattern, EndPattern gg.Pattern
}

/* Reset the current step counter */
func (options *SteppedDrawingOptions) resetStepCount() {
	options.currentStep = 0
}

/* Increment the step counter */
func (options *SteppedDrawingOptions) incStepCount() {
	options.currentStep++
}

/*
 * Implement the Pattern interface !
 * Calls the two underlying pattern implementation, then lerps colors.
 */
func (options *SteppedDrawingOptions) ColorAt(x, y int) color.Color {
	startColor := options.StartPattern.ColorAt(x, y)
	endColor := options.EndPattern.ColorAt(x, y)

	return utils.ColorLerp(startColor, endColor, float64(options.currentStep)/float64(options.Steps-1))
}
