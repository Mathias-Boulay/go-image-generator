package shapes

import "github.com/fogleman/gg"

/* Struct representing parameters for the DrawStepped Blob method */
type SteppedDrawingOptions struct {
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
}
